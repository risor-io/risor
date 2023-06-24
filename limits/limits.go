package limits

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const NoLimit = -1

type Limits interface {

	// IOTimeout returns the maximum amount of time to wait for IO operations.
	IOTimeout() time.Duration

	// MaxBodySize returns the maximum allowed size buffer in bytes. This is
	// relevant for any buffered I/O operation, e.g. reading an HTTP request body.
	MaxBufferSize() int64

	// TrackHTTPRequest returns an error if the HTTP request should not
	// be processed due to exceeding a limit.
	TrackHTTPRequest(*http.Request) error

	// TrackHTTPResponse returns an error if the HTTP response should not
	// be processed due to exceeding a body or content length limit.
	TrackHTTPResponse(*http.Response) error

	// TrackCost returns an error if the given incremental processing cost
	// causes the cost limit to be exceeded.
	TrackCost(cost int) error

	// ReadAll reads from the given reader until EOF or a limit is reached.
	// This counts towards the allocation limit.
	ReadAll(reader io.Reader) ([]byte, error)
}

type contextKey string

const limitsKey = contextKey("tamarin:limits")

// WithCodeFunc adds an CodeFunc to the context, which can be used by
// objects to retrieve the active code at runtime
func WithLimits(ctx context.Context, l Limits) context.Context {
	return context.WithValue(ctx, limitsKey, l)
}

// GetLimits returns Tamarin limits associated with the context, if any.
func GetLimits(ctx context.Context) (Limits, bool) {
	l, ok := ctx.Value(limitsKey).(Limits)
	return l, ok
}

// TrackCost increments the processing cost associated with the context by
// the given amount. If the cost limit is exceeded, an error is returned.
func TrackCost(ctx context.Context, cost int) error {
	l, ok := GetLimits(ctx)
	if !ok {
		return LimitsNotFound
	}
	return l.TrackCost(cost)
}

// LimitsError indicates that a limit was exceeded.
type LimitsError struct {
	message string
}

func (e *LimitsError) Error() string {
	return e.message
}

// NewLimitsError returns a new LimitsError with the given message.
func NewLimitsError(message string, args ...interface{}) error {
	return &LimitsError{message: fmt.Sprintf(message, args...)}
}

// LimitsNotFound is a standard error that indicates limits were expected to be
// present in a context, but were not found.
var LimitsNotFound = NewLimitsError("limit error: limits not found in context")

// ReadAll reads from the given reader until EOF or the limit is reached.
// If the given limit is less than zero, the entire reader is read.
func ReadAll(reader io.Reader, limit int64) ([]byte, error) {
	if limit <= NoLimit {
		return io.ReadAll(reader)
	}
	// Bump limit by one so we can detect if the supplied limit was exceeded,
	// in which case we'll return an error.
	limit++
	body, err := io.ReadAll(io.LimitReader(reader, limit))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) >= limit {
		return nil, NewLimitsError("limit error: data size exceeded limit of %d bytes", limit)
	}
	return body, nil
}
