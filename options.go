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
	"fmt"
	"io"
	"time"
)

// Option is used to configure each probe attempt.
type Option func(opts *options)

// WithKind specifies the kind of probe, i.e. {read,write} {IOPS,bandwidth}.
func WithKind(kind Kind) Option {
	return func(opts *options) {
		opts.Kind = kind
	}
}

// WithDirectory configures the probe to make use of the given directory. The
// underlying volume is what ends up getting measured.
func WithDirectory(dir string) Option {
	return func(opts *options) {
		opts.Directory = dir
	}
}

// WithDuration controls how long we record measurements for.
func WithDuration(dur time.Duration) Option {
	return func(opts *options) {
		opts.Duration = dur
	}
}

// WithRamp controls the ramp-up period before the actual probe attempt.
func WithRamp(ramp time.Duration) Option {
	return func(opts *options) {
		opts.Ramp = ramp
	}
}

// WithSize controls how many bytes are written to disk during the probe.
func WithSize(size uint64) Option {
	return func(opts *options) {
		opts.Size = size
	}
}

// WithMaxRate limits the probe to a maximum bandwidth (if a {read,write}
// bandwidth probe) or IOPS (if a {read,write} IOPS probe).
func WithMaxRate(rate uint64) Option {
	return func(opts *options) {
		opts.MaxRate = rate
	}
}

// WithLoggingTo instructs the liveness module to log to the given io.Writer.
func WithLoggingTo(w io.Writer) Option {
	return func(opts *options) {
		opts.LoggingTo = w
	}
}

type options struct {
	Directory string
	Duration  time.Duration
	Ramp      time.Duration
	Size      uint64
	Kind      Kind
	MaxRate   uint64
	LoggingTo io.Writer
}

func (o *options) validate() error {
	if o.Kind == "" {
		return fmt.Errorf("probe kind unspecified")
	}
	return nil
}
