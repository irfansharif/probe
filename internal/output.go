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

package internal

// Output represents the top-level JSON output for fio.
type Output struct {
	Jobs []Job `json:"jobs"`
}

// Job represents the JSON output for each job.
type Job struct {
	Read  ReadWriteStats `json:"read"`
	Write ReadWriteStats `json:"write"`
}

// ReadWriteStats represents the JSON output for read/write statistics.
type ReadWriteStats struct {
	BWBytes int     `json:"bw_bytes"`
	IOPS    float64 `json:"iops"`

	BWMin     int     `json:"bw_min"`
	BWMax     int     `json:"bw_max"`
	BWAgg     float64 `json:"bw_agg"`
	BWMean    float64 `json:"bw_mean"`
	BWDev     float64 `json:"bw_dev"`
	BWSamples int     `json:"bw_samples"`

	IOPSMin     int     `json:"iops_min"`
	IOPSMax     int     `json:"iops_max"`
	IOPSMean    float64 `json:"iops_mean"`
	IOPSStddev  float64 `json:"iops_stddev"`
	IOPSSamples int     `json:"iops_samples"`
}

// NB: The full JSON output for fio looks as such:
//
// {
//   "fio version" : "fio-3.30",
//   "timestamp" : 1690585048,
//   "timestamp_ms" : 1690585048882,
//   "time" : "Fri Jul 28 18:57:28 2023",
//   "jobs" : [
//     {
//       "jobname" : "write_throughput",
//       "groupid" : 0,
//       "error" : 0,
//       "eta" : 0,
//       "elapsed" : 13,
//       "job options" : {
//         "name" : "write_throughput",
//         "directory" : "dir",
//         "numjobs" : "8",
//         "size" : "11GB",
//         "runtime" : "10s",
//         "ramp_time" : "2s",
//         "ioengine" : "posixaio",
//         "direct" : "1",
//         "verify" : "0",
//         "bs" : "1.0MB",
//         "iodepth" : "64",
//         "rw" : "write",
//         "group_reporting" : "1"
//       },
//       "read" : {
//         "io_bytes" : 0,
//         "io_kbytes" : 0,
//         "bw_bytes" : 0,
//         "bw" : 0,
//         "iops" : 0.000000,
//         "runtime" : 0,
//         "total_ios" : 0,
//         "short_ios" : 0,
//         "drop_ios" : 0,
//         "slat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         },
//         "clat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         },
//         "lat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         },
//         "bw_min" : 0,
//         "bw_max" : 0,
//         "bw_agg" : 0.000000,
//         "bw_mean" : 0.000000,
//         "bw_dev" : 0.000000,
//         "bw_samples" : 0,
//         "iops_min" : 0,
//         "iops_max" : 0,
//         "iops_mean" : 0.000000,
//         "iops_stddev" : 0.000000,
//         "iops_samples" : 0
//       },
//       "write" : {
//         "io_bytes" : 3752351,
//         "io_kbytes" : 3664,
//         "bw_bytes" : 375197,
//         "bw" : 366,
//         "iops" : 375191.380862,
//         "runtime" : 10001,
//         "total_ios" : 3752289,
//         "short_ios" : 0,
//         "drop_ios" : 0,
//         "slat_ns" : {
//           "min" : 0,
//           "max" : 42889000,
//           "mean" : 8255.047179,
//           "stddev" : 329922.495494,
//           "N" : 3752294
//         },
//         "clat_ns" : {
//           "min" : 1000,
//           "max" : 48451000,
//           "mean" : 210681.698074,
//           "stddev" : 1392251.256053,
//           "N" : 3752321,
//           "percentile" : {
//             "1.000000" : 11968,
//             "5.000000" : 20096,
//             "10.000000" : 28032,
//             "20.000000" : 42240,
//             "30.000000" : 57088,
//             "40.000000" : 75264,
//             "50.000000" : 105984,
//             "60.000000" : 179200,
//             "70.000000" : 211968,
//             "80.000000" : 222208,
//             "90.000000" : 240640,
//             "95.000000" : 337920,
//             "99.000000" : 806912,
//             "99.500000" : 1351680,
//             "99.900000" : 32112640,
//             "99.950000" : 34340864,
//             "99.990000" : 37486592
//           }
//         },
//         "lat_ns" : {
//           "min" : 6000,
//           "max" : 48452000,
//           "mean" : 218936.727695,
//           "stddev" : 1430845.458417,
//           "N" : 3752321
//         },
//         "bw_min" : 290,
//         "bw_max" : 435,
//         "bw_agg" : 98.798231,
//         "bw_mean" : 362.157895,
//         "bw_dev" : 4.330379,
//         "bw_samples" : 152,
//         "iops_min" : 301919,
//         "iops_max" : 448494,
//         "iops_mean" : 374900.210526,
//         "iops_stddev" : 4417.469684,
//         "iops_samples" : 152
//       },
//       "trim" : {
//         "io_bytes" : 0,
//         "io_kbytes" : 0,
//         "bw_bytes" : 0,
//         "bw" : 0,
//         "iops" : 0.000000,
//         "runtime" : 0,
//         "total_ios" : 0,
//         "short_ios" : 0,
//         "drop_ios" : 0,
//         "slat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         },
//         "clat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         },
//         "lat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         },
//         "bw_min" : 0,
//         "bw_max" : 0,
//         "bw_agg" : 0.000000,
//         "bw_mean" : 0.000000,
//         "bw_dev" : 0.000000,
//         "bw_samples" : 0,
//         "iops_min" : 0,
//         "iops_max" : 0,
//         "iops_mean" : 0.000000,
//         "iops_stddev" : 0.000000,
//         "iops_samples" : 0
//       },
//       "sync" : {
//         "total_ios" : 0,
//         "lat_ns" : {
//           "min" : 0,
//           "max" : 0,
//           "mean" : 0.000000,
//           "stddev" : 0.000000,
//           "N" : 0
//         }
//       },
//       "job_runtime" : 80000,
//       "usr_cpu" : 5.226250,
//       "sys_cpu" : 10.048750,
//       "ctx" : 2390859,
//       "majf" : 0,
//       "minf" : 81,
//       "iodepth_level" : {
//         "1" : 1.597745,
//         "2" : 8.348291,
//         "4" : 21.887200,
//         "8" : 61.037383,
//         "16" : 7.129382,
//         "32" : 0.000000,
//         ">=64" : 0.000000
//       },
//       "iodepth_submit" : {
//         "0" : 0.000000,
//         "4" : 100.000000,
//         "8" : 0.000000,
//         "16" : 0.000000,
//         "32" : 0.000000,
//         "64" : 0.000000,
//         ">=64" : 0.000000
//       },
//       "iodepth_complete" : {
//         "0" : 0.000000,
//         "4" : 94.514181,
//         "8" : 2.073844,
//         "16" : 3.411975,
//         "32" : 0.000000,
//         "64" : 0.000000,
//         ">=64" : 0.000000
//       },
//       "latency_ns" : {
//         "2" : 0.000000,
//         "4" : 0.000000,
//         "10" : 0.000000,
//         "20" : 0.000000,
//         "50" : 0.000000,
//         "100" : 0.000000,
//         "250" : 0.000000,
//         "500" : 0.000000,
//         "750" : 0.000000,
//         "1000" : 0.000000
//       },
//       "latency_us" : {
//         "2" : 0.010000,
//         "4" : 0.010000,
//         "10" : 0.425287,
//         "20" : 4.188936,
//         "50" : 20.835282,
//         "100" : 23.054701,
//         "250" : 42.379865,
//         "500" : 7.327607,
//         "750" : 0.693017,
//         "1000" : 0.352798
//       },
//       "latency_ms" : {
//         "2" : 0.386564,
//         "4" : 0.074701,
//         "10" : 0.055859,
//         "20" : 0.046665,
//         "50" : 0.172242,
//         "100" : 0.000000,
//         "250" : 0.000000,
//         "500" : 0.000000,
//         "750" : 0.000000,
//         "1000" : 0.000000,
//         "2000" : 0.000000,
//         ">=2000" : 0.000000
//       },
//       "latency_depth" : 64,
//       "latency_target" : 0,
//       "latency_percentile" : 100.000000,
//       "latency_window" : 0
//     }
//   ]
// }
