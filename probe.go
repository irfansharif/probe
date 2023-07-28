// Copyright 2023 Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/irfansharif/probe/internal"
	"github.com/shirou/gopsutil/v3/disk"
)

// Kind of probe; one of {read,write} {bandwidth,IOPS}.
type Kind string

const (
	ReadBandwidth  Kind = "read_bandwidth"
	WriteBandwidth Kind = "write_bandwidth"
	ReadIOPS       Kind = "read_iops"
	WriteIOPS      Kind = "write_iops"
)

// Supported returns whether the probe library is supported (thin check for
// whether 'fio' is installed and accessible).
func Supported() bool {
	_, err := exec.LookPath("fio")
	return err == nil
}

// TODO: Write controller that's able to observe current rate (for CRDB process,
//       or aggregate bandwidth, or compaction bandwidth) and probe higher to
//       see if we're able to push out higher. I think we'd need to observe
//       aggregate bandwidth. Probing up by 50MiB/s once observed bandwidth
//       pre-probe is stable (low stdev), to see if aggregate bandwidth does
//       increase. And probe down if IO latencies are unacceptable. Looking at
//       either PSI metrics, or something else.

// Probe disks for their capacity, i.e. {read,write} {bandwidth,IOPS}.
func Probe(ctx context.Context, opts ...Option) (_ uint64, err error) {
	// Test {read,write} throughput by performing sequential {read,writes} with
	// multiple parallel streams (8+), using an I/O block size of 1 MB and an
	// I/O depth of at least 64.
	//
	// 	fio --name=write_throughput --directory=$TEST_DIR --numjobs=8 \
	//	  --size=10G --time_based --runtime=60s --ramp_time=2s --ioengine=libaio \
	//	  --direct=1 --verify=0 --bs=1M --iodepth=64 --rw=write \
	//	  --group_reporting=1
	//
	//	fio --name=read_throughput --directory=$TEST_DIR --numjobs=8 \
	//	  --size=10G --time_based --runtime=60s --ramp_time=2s --ioengine=libaio \
	//	  --direct=1 --verify=0 --bs=1M --iodepth=64 --rw=read \
	//	  --group_reporting=1
	//
	// Test {read,write} IOPS by performing random {read,writes}, using an I/O
	// block size of 4 KB and an I/O depth of at least 64.
	//
	//	fio --name=write_iops --directory=$TEST_DIR
	//    --size=10G --time_based --runtime=60s --ramp_time=2s --ioengine=libaio \
	//    --direct=1 --verify=0 --bs=4K --iodepth=64 --rw=randwrite \
	//    --group_reporting=1
	//
	//	fio --name=read_iops --directory=$TEST_DIR \
	//    --size=10G --time_based --runtime=60s --ramp_time=2s --ioengine=libaio \
	//    --direct=1 --verify=0 --bs=4K --iodepth=64 --rw=randread \
	//    --group_reporting=1

	o := &options{
		Duration:  60 * time.Second,
		Ramp:      2 * time.Second,
		Size:      10 << 30, // 10 GiB
		LoggingTo: io.Discard,
	}
	for _, opt := range opts {
		opt(o)
	}
	if err := o.validate(); err != nil {
		return 0, err
	}

	if err := os.RemoveAll(o.Directory); err != nil {
		// Nuke left-over state, if any. We don't want to accrete storage use
		// across {failed,} runs.
		return 0, err
	}
	if err := os.MkdirAll(o.Directory, 0755); err != nil {
		return 0, err
	}
	defer func() {
		if err2 := os.RemoveAll(o.Directory); err2 != nil {
			if err == nil {
				err = err2
			} else {
				err = fmt.Errorf("%s: %s", err2, err)
			}
		}
	}()

	var ioengine = "libaio"
	if runtime.GOOS == "darwin" {
		ioengine = "posixaio"
	}

	var args []string
	args = append(args,
		"--name", string(o.Kind),
		"--directory", o.Directory,
		"--size", fmt.Sprint(o.Size),
		"--time_based", "--runtime", fmt.Sprintf("%ds", int(o.Duration.Seconds())),
		"--ramp_time", fmt.Sprintf("%ds", int(o.Ramp.Seconds())),
		"--ioengine", ioengine,
		"--direct", "1",
		"--verify", "0",
		"--bs", fmt.Sprint(1<<20), /* 1MiB */
		"--iodepth", "64",
		"--group_reporting=1",
		"--output-format", "json",
	)

	if (o.Kind == ReadBandwidth) || (o.Kind == WriteBandwidth) {
		args = append(args, "--numjobs", "8")
		// Limit aggregate disk use across jobs.
		o.Size /= 8
	}

	switch o.Kind {
	case ReadBandwidth:
		args = append(args, "--rw", "read")
	case WriteBandwidth:
		args = append(args, "--rw", "write")
	case ReadIOPS:
		args = append(args, "--rw", "randread")
	case WriteIOPS:
		args = append(args, "--rw", "randwrite")
	default:
		return 0, fmt.Errorf("invalid kind: %s", o.Kind)
	}

	if (o.Kind == ReadBandwidth) || (o.Kind == WriteBandwidth) {
		// Use 1MiB block sizes for bandwidth probes.
		args = append(args, "--bs", fmt.Sprint(1<<20) /* 1MiB */)
	} else {
		// Use 4KiB block sizes for IOPS probes.
		args = append(args, "--bs", fmt.Sprint(4<<10) /* 4KiB */)
	}

	if o.MaxRate != 0 {
		if (o.Kind == ReadBandwidth) || (o.Kind == WriteBandwidth) {
			// We want to preserve a max rate across 8 jobs, so divide
			// accordingly.
			args = append(args, "--rate", fmt.Sprint(o.MaxRate/8))
		} else {
			args = append(args, "--rate_iops", fmt.Sprint(o.MaxRate))
		}
	}

	usage, err := disk.Usage(o.Directory)
	if err != nil {
		return 0, err
	}
	if limit := o.Size + (5 << 30); usage.Free < limit {
		return 0, fmt.Errorf("insufficient disk space: %s, want %s",
			humanize.IBytes(usage.Free),
			humanize.IBytes(limit))
	}

	cmd := exec.CommandContext(ctx, "fio", args...)

	if false {
		// Sometimes useful for debugging.
		fmt.Println(cmd.String())
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = o.LoggingTo.Write(output)
		return 0, err
	}

	var fiout internal.Output
	if err := json.Unmarshal(output, &fiout); err != nil {
		return 0, err
	}

	switch o.Kind {
	case ReadBandwidth:
		return uint64(fiout.Jobs[0].Read.BWBytes), nil
	case WriteBandwidth:
		return uint64(fiout.Jobs[0].Write.BWBytes), nil
	case ReadIOPS:
		return uint64(fiout.Jobs[0].Read.IOPS), nil
	case WriteIOPS:
		return uint64(fiout.Jobs[0].Write.IOPS), nil
	default:
		return 0, fmt.Errorf("invalid kind: %s", o.Kind)
	}
}
