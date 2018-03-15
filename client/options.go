package client

import (
	"io"
	"os"
	"time"
)

// Options represents the options used during the creation a new client instance
type Options struct {
	Hooks                 Hooks
	DefaultRequestTimeout time.Duration
	WarnLog               io.Writer
	ErrorLog              io.Writer
}

// SetDefaults sets default values for undefined required options
func (opts *Options) SetDefaults() {
	opts.Hooks.SetDefaults()

	if opts.DefaultRequestTimeout < 1 {
		opts.DefaultRequestTimeout = 60 * time.Second
	}

	if opts.WarnLog == nil {
		opts.WarnLog = os.Stdout
	}

	if opts.ErrorLog == nil {
		opts.ErrorLog = os.Stderr
	}
}
