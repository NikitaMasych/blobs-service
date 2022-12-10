package ctx

import (
	"context"
	"gitlab.com/distributed_lab/kit/pgdb"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logKey ctxKey = iota
	dbKey
)

func SetLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logKey).(*logan.Entry)
}

func SetDB(db *pgdb.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbKey, db)
	}
}

func DB(r *http.Request) *pgdb.DB {
	return r.Context().Value(dbKey).(*pgdb.DB)
}
