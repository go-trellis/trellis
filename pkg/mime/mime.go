package mime

const (
	ContentTypeJson    = "application/json"
	ContentTypeJsonBom = "application/json; charset=UTF-8"
)

const (
	HeaderKeyTraceID          = "X-Trace-Id"
	HeaderKeyClientIP         = "X-Client-Ip"
	HeaderKeyRequestIP        = "X-Request-Ip"
	HeaderKeyContentType      = "Content-Type"
	HeaderKeyUserAgent        = "User-Agent"
	HeaderKeyRequestURIPath   = "Request-URI-Path"
	HeaderKeyRequestURIQuery  = "Request-URI-Query"
	HeaderKeyRequestURIMethod = "Request-URI-Method"
)

type TraceInfo struct {
	TraceID         string `yaml:"X-Trace-Id" json:"X-Trace-Id"`
	ClientIP        string `yaml:"X-Client-Ip" json:"X-Client-Ip"`
	RequestIP       string `yaml:"X-Request-Ip" json:"X-Request-Ip"`
	ContentType     string `yaml:"Content-Type" json:"Content-Type"`
	UserAgent       string `yaml:"User-Agent" json:"User-Agent"`
	RequestURIPath  string `yaml:"Request-URI-Path" json:"Request-URI-Path"`
	RequestURIQuery string `yaml:"Request-URI-Query" json:"Request-URI-Query"`
}
