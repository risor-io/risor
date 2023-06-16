package limits

type HTTP struct {
	MaxBodyLength    int64 `json:"max_body_length"`
	MaxContentLength int64 `json:"max_content_length"`
	Timeout          int64 `json:"timeout"`
}

type Limits struct {
	HTTP `json:"http"`
}
