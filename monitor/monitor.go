package monitor

import (
	"context"
)

type Options struct {
	context     context.Context
	maxAttempts uint
}

type OptFn func(*Options)

func WithMaxAttempts(attempts uint) OptFn {
	return func(config *Options) {
		config.maxAttempts = attempts
	}
}

func WithContext(ctx context.Context) OptFn {
	return func(options *Options) {
		options.context = ctx
	}
}

func defaultOpts() Options {
	return Options{
		context:     context.Background(),
		maxAttempts: uint(2),
	}
}

func Do(fn func() error, optFns ...OptFn) error {
	var (
		err error
		n   uint

		canRetry = true
		opts     = defaultOpts()
		attempts = opts.maxAttempts
	)

	for _, optFn := range optFns {
		optFn(&opts)
	}

	for canRetry {
		if err = fn(); err == nil {
			return nil
		}

		n++
		canRetry = n < attempts
	}

	return err
}
