package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type LobCloser struct {
	tx pgx.Tx
	*pgx.LargeObject
}

func (l *LobCloser) Close() error {
	defer l.tx.Rollback(context.Background())
	return l.LargeObject.Close()
}
