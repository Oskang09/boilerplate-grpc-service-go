package api

import (
	"context"
	"encoding/json"
	"runtime"
	"service/app/bootstrap"
	"service/app/config"
	"service/app/repository"
	"service/package/grpc/errors"
	"service/package/redis"
	"service/package/redsync"
	v1 "service/protobuf/go/v1"

	"github.com/RevenueMonster/sqlike/sqlike"
	"github.com/opentracing/opentracing-go"

	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

type example struct {
	database   *sqlike.Database
	repository *repository.Repository
	redsync    *redsync.Redsync
	redis      *redis.Client
	validator  *validator.Validate

	v1.UnimplementedExampleServiceServer
}

func (*example) bind(request interface{}, i interface{}) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	return nil
}
func (e *example) validate(i interface{}) error {
	return e.validate(i)
}

func New(bs *bootstrap.Bootstrap) *example {
	handler := new(example)
	handler.database = bs.Database
	handler.repository = bs.Repository
	handler.redsync = bs.Redsync
	handler.redis = bs.Redis

	validator := validator.New()
	en_translations.RegisterDefaultTranslations(validator, config.ValidatorTranslator)
	handler.validator = validator

	return handler
}

func responseError(ctx context.Context, statusCode int32, code string, msg interface{}) error {
	val := ""
	switch v := msg.(type) {
	case nil:
	case error:
		val = v.Error()
	case string:
		val = v
	default:
		data, _ := json.Marshal(v)
		val = string(data)
	}

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("grpc.status", statusCode)
		span.SetTag("grpc.code", code)

		_, file, no, ok := runtime.Caller(1)
		if ok {
			span.LogKV(
				"file", file,
				"line", no,
				"erorr", val,
			)
		}
	}

	return errors.Service(statusCode, code, val)
}
