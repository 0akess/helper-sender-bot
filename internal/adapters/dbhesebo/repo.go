package dbhesebo

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pgx struct {
	Db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewStorage(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		Db: pool,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
