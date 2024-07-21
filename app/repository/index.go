package repository

import (
	"service/app/model"

	"github.com/RevenueMonster/sqlike/sqlike"
)

type tableContext interface {
	Table(string) *sqlike.Table
}

type Repository struct {
	Example Example
}

func New(db tableContext) *Repository {
	return &Repository{
		Example{newRepository[model.Example](db)},
	}
}
