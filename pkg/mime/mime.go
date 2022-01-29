package mime

const (
	ContentTypeJson = "application/json"
)

const (
	HeaderKeyTraceID     = "X-TRACE-ID"
	HeaderKeyClientIP    = "X-CLIENT-IP"
	HeaderKeyRequestIP   = "X-REQUEST-IP"
	HeaderKeyRequestID   = "X-REQUEST-ID"
	HeaderKeyContentType = "Content-Type"
)

type TraceInfo struct {
	TraceID   string `yaml:"X-TRACE-ID" json:"X-TRACE-ID"`
	RequestID string `yaml:"X-REQUEST-ID" json:"X-REQUEST-ID"`
	ClientIP  string `yaml:"X-CLIENT-IP" json:"X-CLIENT-IP"`
	RequestIP string `yaml:"X-REQUEST-IP" json:"X-REQUEST-IP"`
}
