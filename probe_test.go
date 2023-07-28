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

package probe_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/irfansharif/probe"
)

func TestSupported(t *testing.T) {
	t.Logf("supported = %t", probe.Supported())
}

var logger = log.New(os.Stdout, "[probe] ", log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)

var opts = []probe.Option{
	probe.WithDirectory("dir"),
	probe.WithDuration(10 * time.Second),
	probe.WithRamp(2 * time.Second),
	probe.WithSize(5 << 30 /* 5 GiB */),
	probe.WithLoggingTo(logger.Writer()),
}

func TestWriteIOPS(t *testing.T) {
	for _, withMax := range []bool{true, false} {
		t.Run(fmt.Sprintf("with-max=%t", withMax), func(t *testing.T) {
			ctx := context.Background()
			opts := append(opts, probe.WithKind(probe.WriteIOPS))
			if withMax {
				opts = append(opts, probe.WithMaxRate(1000))
			}
			iops, err := probe.Probe(ctx, opts...)
			if err != nil {
				t.Fatal(err)
			}

			withMaxInfix := ""
			if withMax {
				withMaxInfix = fmt.Sprintf("(max = %d) ", 1000)
			}
			t.Logf("write iops %s= %d", withMaxInfix, iops)
		})
	}
}

func TestWriteBandwidth(t *testing.T) {
	for _, withMax := range []bool{true, false} {
		t.Run(fmt.Sprintf("with-max=%t", withMax), func(t *testing.T) {
			ctx := context.Background()
			opts := append(opts, probe.WithKind(probe.WriteBandwidth))
			if withMax {
				opts = append(opts, probe.WithMaxRate(80<<20))
			}
			bw, err := probe.Probe(ctx, opts...)
			if err != nil {
				t.Fatal(err)
			}

			withMaxInfix := ""
			if withMax {
				withMaxInfix = fmt.Sprintf("(max = %s/s) ", humanize.IBytes(80<<20))
			}
			t.Logf("write bandwidth %s= %s/s", withMaxInfix, humanize.IBytes(bw))
		})
	}
}

func TestReadIOPS(t *testing.T) {
	for _, withMax := range []bool{true, false} {
		t.Run(fmt.Sprintf("with-max=%t", withMax), func(t *testing.T) {
			ctx := context.Background()
			opts := append(opts, probe.WithKind(probe.ReadIOPS))
			if withMax {
				opts = append(opts, probe.WithMaxRate(1000))
			}
			iops, err := probe.Probe(ctx, opts...)
			if err != nil {
				t.Fatal(err)
			}

			withMaxInfix := ""
			if withMax {
				withMaxInfix = fmt.Sprintf("(max = %d) ", 1000)
			}
			t.Logf("read iops %s= %d", withMaxInfix, iops)
		})
	}
}

func TestReadBandwidth(t *testing.T) {
	for _, withMax := range []bool{true, false} {
		t.Run(fmt.Sprintf("with-max=%t", withMax), func(t *testing.T) {
			ctx := context.Background()
			opts := append(opts, probe.WithKind(probe.ReadBandwidth))
			if withMax {
				opts = append(opts, probe.WithMaxRate(80<<20))
			}
			bw, err := probe.Probe(ctx, opts...)
			if err != nil {
				t.Fatal(err)
			}

			withMaxInfix := ""
			if withMax {
				withMaxInfix = fmt.Sprintf("(max = %s/s) ", humanize.IBytes(80<<20))
			}
			t.Logf("read bandwidth %s= %s/s", withMaxInfix, humanize.IBytes(bw))
		})
	}
}
