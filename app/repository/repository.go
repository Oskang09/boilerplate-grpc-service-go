package repository

import (
	"context"
	"reflect"
	"service/app/model"

	"github.com/RevenueMonster/sqlike/sql/expr"
	"github.com/RevenueMonster/sqlike/sqlike"
	"github.com/RevenueMonster/sqlike/sqlike/actions"
	"github.com/RevenueMonster/sqlike/sqlike/options"
)

type PaginateOptions struct {
	Limit   uint
	Cursor  string
	Queries []interface{}
	Sorts   []interface{}
}

func NewPaginateOptions() *PaginateOptions {
	opt := new(PaginateOptions)
	opt.Limit = 50
	opt.Cursor = ""
	opt.Queries = make([]interface{}, 0)
	return opt
}

func newRepository[T model.TableModel](db tableContext) baseRepository[T] {
	return baseRepository[T]{db.Table((*new(T)).Table())}
}

type baseRepository[T any] struct {
	table *sqlike.Table
}

// ---------
// TODO :-
//
// Build an option functionality to support caching when find by key
// So when upsert / save function run will clear the cache by key
// ---------
func (r baseRepository[T]) Find(ctx context.Context, key string) (*T, error) {
	query := actions.FindOne().Where(expr.Equal("Key", key))
	result := r.table.FindOne(ctx, query)
	model := new(T)
	if err := result.Decode(model); err != nil {
		return nil, err
	}
	return model, nil
}

func (r baseRepository[T]) Create(ctx context.Context, model *T) (err error) {
	_, err = r.table.InsertOne(ctx, model)
	return
}

func (r baseRepository[T]) Migrate(ctx context.Context) (err error) {
	err = r.table.Migrate(ctx, new(T))
	return
}

func (r baseRepository[T]) Upsert(ctx context.Context, model *T) (err error) {
	_, err = r.table.InsertOne(
		ctx,
		model,
		options.InsertOne().SetMode(options.InsertOnDuplicate),
	)
	return
}

func (r baseRepository[T]) Delete(ctx context.Context, model *T) (err error) {
	err = r.table.DestroyOne(ctx, model)
	return
}

func (r baseRepository[T]) Paginate(ctx context.Context, opts *PaginateOptions) ([]*T, string, error) {

	if opts.Limit > 100 {
		opts.Limit = 100
	}

	if opts.Limit <= 0 {
		opts.Limit = 50
	}

	query := actions.Paginate().Limit(opts.Limit + 1).Where(opts.Queries...).OrderBy(opts.Sorts...)
	pg, err := r.table.Paginate(ctx, query, options.Paginate().SetDebug(true))
	if err != nil {
		return nil, "", err
	}

	if opts.Cursor != "" {
		err = pg.NextCursor(ctx, opts.Cursor)
		if err != nil {
			return nil, "", err
		}
	}

	models := make([]*T, 0)
	err = pg.All(models)
	if err != nil {
		return nil, "", err
	}

	ptr := reflect.ValueOf(models).Elem()
	if v := ptr.Len(); v == int(opts.Limit)+1 {
		data := ptr.Index(int(opts.Limit) - 1) // get last value because when get will be (Limit+1)
		if data.Kind() == reflect.Ptr {
			data = data.Elem()
		}

		// remove next cursor value
		ptr.Set(ptr.Slice(0, v-1))

		keyValue := data.FieldByName("Key")
		key := keyValue.Interface()
		return models, key.(string), nil
	}

	return models, "", nil
}
