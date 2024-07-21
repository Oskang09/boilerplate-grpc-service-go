package background

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TaskOption struct {
	ctx          context.Context
	name         string
	timeout      time.Duration
	retryOptions []retry.Option
	logs         []log.Field
}

type TaskOptionFn func(*TaskOption)

type Task func(context.Context, opentracing.Span) error

func WithParentContext(ctx context.Context) TaskOptionFn {
	return func(to *TaskOption) {
		to.ctx = ctx
	}
}

func WithName(name string) TaskOptionFn {
	return func(to *TaskOption) {
		to.name = name
	}
}

func WithTimeout(timeout time.Duration) TaskOptionFn {
	return func(to *TaskOption) {
		to.timeout = timeout
	}
}

func WithLogs(logs ...log.Field) TaskOptionFn {
	return func(to *TaskOption) {
		to.logs = logs
	}
}

func WithRetryOptions(options ...retry.Option) TaskOptionFn {
	return func(to *TaskOption) {
		to.retryOptions = options
	}
}

func RunTask(fn Task, options ...TaskOptionFn) {
	to := new(TaskOption)
	for _, opt := range options {
		opt(to)
	}

	span := opentracing.SpanFromContext(to.ctx)
	ref := make([]opentracing.StartSpanOption, 0)
	if span != nil {
		ref = append(ref, opentracing.ChildOf(span.Context()))
	}

	span = opentracing.StartSpan(to.name, ref...)
	defer span.Finish()

	ctx := context.Background()
	if to.timeout > 0 {
		newCtx, cancel := context.WithTimeout(ctx, to.timeout)
		defer cancel()

		ctx = newCtx
	}

	ctx = opentracing.ContextWithSpan(ctx, span)
	for _, value := range to.logs {
		span.LogFields(value)
	}

	var err error
	if len(to.retryOptions) > 0 {
		err = retry.Do(
			func() error {
				return fn(ctx, span)
			},
			to.retryOptions...,
		)
	} else {
		err = fn(ctx, span)
	}

	if err != nil {
		span.LogFields(log.Error(err))
	}
}
