package util

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

const (
	RequestID = contextKey("RequestID")
	LayoutISO = "2006-01-02 15:04:05"
)

type contextKey string

func GetRequestID(ctx context.Context) (reqID string) {
	if ctx == nil {
		return
	}
	reqID, _ = ctx.Value(RequestID).(string)
	return
}

func SetRequestID(ctx context.Context) context.Context {
	if len(GetRequestID(ctx)) == 0 {
		reqid := uuid.NewV4()
		return context.WithValue(ctx, RequestID, reqid.String())
	}
	return ctx
}
