package redis

import (
	"context"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
)

// Setup :
type Setup struct {
	RedisHost      string
	RedisPassword  string
	OpenTracingLog bool
}

// Client :
type Client struct {
	Client         *redis.Client
	openTracingLog bool
}

// Config :
func Config(rs *Setup) *Client {
	redisHost := rs.RedisHost
	redisPassword := rs.RedisPassword

	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	setup := &Client{
		Client:         client,
		openTracingLog: rs.OpenTracingLog,
	}

	return setup
}

// SetKeyWithExpired :
func (r *Client) SetKeyWithExpired(ctx context.Context, key string, payload interface{}, expiredTimeInSecond int) error {

	if r.openTracingLog {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			newSpan := opentracing.StartSpan("redis", opentracing.ChildOf(span.Context()))
			defer newSpan.Finish()

			_, file, line, ok := runtime.Caller(1)
			if ok {
				newSpan.LogKV(
					"action", "SetKeyWithExpired",
					"expiresInSecond", expiredTimeInSecond,
					"key", key,
					"payload", payload,
					"file", file,
					"line", line,
				)
			}
		}
	}

	return r.Client.Set(ctx, key, payload, time.Duration(expiredTimeInSecond)*time.Second).Err()
}

// GetPayloadByKey :
func (r *Client) GetPayloadByKey(ctx context.Context, key string) (string, error) {

	data, err := r.Client.Get(ctx, key).Result()

	if r.openTracingLog {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			newSpan := opentracing.StartSpan("redis", opentracing.ChildOf(span.Context()))
			defer newSpan.Finish()

			_, file, line, ok := runtime.Caller(1)
			if ok {
				newSpan.LogKV(
					"action", "GetPayloadByKey",
					"key", key,
					"data", data,
					"file", file,
					"line", line,
				)
			}
		}
	}

	return data, err
}

// Delete :
func (r *Client) Delete(ctx context.Context, key ...string) error {

	if r.openTracingLog {
		span := opentracing.SpanFromContext(ctx)
		if span != nil {
			newSpan := opentracing.StartSpan("redis", opentracing.ChildOf(span.Context()))
			defer newSpan.Finish()

			_, file, line, ok := runtime.Caller(1)
			if ok {
				newSpan.LogKV(
					"action", "Delete",
					"key", key,
					"file", file,
					"line", line,
				)
			}
		}
	}

	return r.Client.Del(ctx, key...).Err()
}
