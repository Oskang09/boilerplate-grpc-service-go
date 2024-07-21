package redsync

import (
	"context"
	"runtime"

	"github.com/go-redsync/redsync/v4"
	"github.com/opentracing/opentracing-go"
)

type Redsync struct {
	client         *redsync.Redsync
	opentracingLog bool
}

func New(opentracingLog bool, client *redsync.Redsync) *Redsync {
	return &Redsync{client, opentracingLog}
}

// NewMutex returns a new distributed mutex with given name.
func (r *Redsync) NewMutex(ctx context.Context, name string, options ...redsync.Option) *redsync.Mutex {

	if r.opentracingLog {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			newSpan := opentracing.StartSpan("redsync-mutex", opentracing.ChildOf(span.Context()))
			defer newSpan.Finish()

			_, file, line, ok := runtime.Caller(1)
			if ok {
				newSpan.LogKV(
					"name", name,
					"file", file,
					"line", line,
				)
			}
		}
	}

	return r.client.NewMutex(name, options...)
}

type Lockable interface {
	LockerID() string
}

func (r *Redsync) NewModelMutex(ctx context.Context, entity Lockable, options ...redsync.Option) *redsync.Mutex {

	if r.opentracingLog {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			newSpan := opentracing.StartSpan("redsync-model-mutex", opentracing.ChildOf(span.Context()))
			defer newSpan.Finish()

			_, file, line, ok := runtime.Caller(1)
			if ok {
				newSpan.LogKV(
					"name", entity.LockerID(),
					"file", file,
					"line", line,
				)
			}
		}
	}

	return r.client.NewMutex(entity.LockerID(), options...)
}
