package ctx

import (
	"blobs/internal/data"
	"context"
	"gitlab.com/tokend/connectors/submit"
	"gitlab.com/tokend/go/xdrbuild"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logKey ctxKey = iota
	blobQKey
	submitterKey
	builderKey
)

func SetLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logKey).(*logan.Entry)
}

func SetBlobQ(q data.Blobs) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blobQKey, q)
	}
}

func BlobQ(r *http.Request) data.Blobs {
	return r.Context().Value(blobQKey).(data.Blobs).New()
}

func SetSubmitter(d *submit.Submitter) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, submitterKey, d)
	}
}

func Submitter(r *http.Request) *submit.Submitter {
	return r.Context().Value(submitterKey).(*submit.Submitter)
}

func SetBuilder(b *xdrbuild.Builder) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, builderKey, b)
	}
}

func Builder(r *http.Request) *xdrbuild.Builder {
	return r.Context().Value(builderKey).(*xdrbuild.Builder)
}
