package limits

import (
	"io"
	"net/http"
	"sync"
	"time"
)

type StandardLimits struct {
	// Configuration
	ioTimeout           time.Duration
	maxBufferSize       int64
	maxHttpRequestCount int64
	maxCost             int64
	// Metrics
	httpRequestsCount int64
	cost              int64
	// Thread safety
	mutex sync.Mutex
}

func (l *StandardLimits) IOTimeout() time.Duration {
	return l.ioTimeout
}

func (l *StandardLimits) MaxBufferSize() int64 {
	return l.maxBufferSize
}

func (l *StandardLimits) TrackHTTPRequest(req *http.Request) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.httpRequestsCount++
	if l.maxHttpRequestCount > NoLimit && l.httpRequestsCount > l.maxHttpRequestCount {
		return NewLimitsError("limit error: reached maximum number of http requests (%d)", l.maxHttpRequestCount)
	}
	return nil
}

func (l *StandardLimits) TrackHTTPResponse(resp *http.Response) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.maxBufferSize > NoLimit && resp.ContentLength >= 0 && resp.ContentLength > l.maxBufferSize {
		return NewLimitsError("limit error: http response content length exceeds maximum allowed buffer size of %d bytes (got %d bytes)",
			l.maxBufferSize, resp.ContentLength)
	}
	return nil
}

func (l *StandardLimits) TrackCost(cost int) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.cost += int64(cost)
	if l.maxCost > NoLimit && l.cost > l.maxCost {
		return NewLimitsError("limit error: reached maximum processing cost (%d)", l.maxCost)
	}
	return nil
}

func (l *StandardLimits) ReadAll(reader io.Reader) ([]byte, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if l.maxBufferSize <= 0 {
		return io.ReadAll(reader)
	}
	bytes, err := ReadAll(reader, l.maxBufferSize)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Option is a function that configures a Limits instance.
type Option func(*StandardLimits)

// WithIOTimeout sets the maximum amount of time to wait for IO operations.
func WithIOTimeout(timeout time.Duration) Option {
	return func(l *StandardLimits) {
		l.ioTimeout = timeout
	}
}

// WithMaxBufferSize sets the maximum allowed size buffer in bytes. This is
// relevant for any buffered I/O operation, e.g. reading an HTTP request body.
func WithMaxBufferSize(size int64) Option {
	return func(l *StandardLimits) {
		l.maxBufferSize = size
	}
}

// WithMaxHttpRequestCount sets the maximum number of HTTP requests that
// are allowed to be processed.
func WithMaxHttpRequestCount(count int64) Option {
	return func(l *StandardLimits) {
		l.maxHttpRequestCount = count
	}
}

// WithMaxCost sets the maximum number allowed processing cost.
func WithMaxCost(cost int64) Option {
	return func(l *StandardLimits) {
		l.maxCost = cost
	}
}

// New creates a new Limits instance with the given options.
func New(opts ...Option) Limits {
	l := &StandardLimits{
		maxBufferSize:       NoLimit,
		maxHttpRequestCount: NoLimit,
		maxCost:             NoLimit,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}
