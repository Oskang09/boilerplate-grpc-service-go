package model

import (
	"time"

	"github.com/RevenueMonster/sqlike/types"
)

type Example struct {
	Key             *types.Key
	Name            string
	CreatedDateTime time.Time
	UpdatedDateTime time.Time
}

func (Example) Table() string {
	return "Example"
}
