/*
Copyright © 2022 Henry Huang <hhh@rutcode.com>
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
