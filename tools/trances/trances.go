package trances

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

func NewTraceID() string {
	return xid.New().String()
}

func NewTraceIDToRequest(r *http.Request, traceID string) *http.Request {
	ctx := context.WithValue(r.Context(), "trances_id", traceID)
	return r.WithContext(ctx)
}
