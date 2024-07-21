package repository

import (
	"context"
	"service/app/model"

	"github.com/RevenueMonster/sqlike/sql/expr"
	"github.com/RevenueMonster/sqlike/sqlike/actions"
)

type Example struct {
	baseRepository[model.Example]
}

func (ex Example) FindByName(ctx context.Context, name string) (*model.Example, error) {
	result := ex.table.FindOne(ctx, actions.FindOne().Where(
		expr.Equal("Name", name),
	))

	example := new(model.Example)
	if err := result.Decode(example); err != nil {
		return nil, err
	}
	return example, nil
}
