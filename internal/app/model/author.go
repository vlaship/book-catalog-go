package model

import (
	"github.com/vlaship/book-catalog-go/internal/app/types"
	"time"
)

type Author struct {
	ID   types.ID  `db:"id"`
	Name string    `db:"name"`
	Dob  time.Time `db:"dob"`
}
