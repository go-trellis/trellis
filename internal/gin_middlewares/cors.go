/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

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

package gin_middlewares

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"trellis.tech/go-trellis/common.v0/config"
	"trellis.tech/trellis.v0/service"
)

func LoadCors(conf config.Config) gin.HandlerFunc {

	var corsConf cors.Config
	if conf == nil {
		corsConf = cors.DefaultConfig()
		corsConf.AllowMethods = []string{"POST"}
		corsConf.AllowOrigins = []string{"*"}
		corsConf.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		corsConf = cors.Config{
			AllowOrigins:     conf.GetStringList("allow-origins"),
			AllowMethods:     conf.GetStringList("allow-methods"),
			AllowHeaders:     conf.GetStringList("allow-headers"),
			ExposeHeaders:    conf.GetStringList("expose-headers"),
			AllowCredentials: conf.GetBoolean("allow-credentials", false),
			MaxAge:           conf.GetTimeDuration("max-age", time.Hour*12),
		}

		corsConf.AllowOriginFunc = wildcardMatchFunc(corsConf.AllowOrigins)
	}

	corsConf.AllowHeaders = append(corsConf.AllowHeaders,
		service.HeaderXAPI, service.HeaderXClientIP, service.HeaderOrigin,
		service.HeaderContentLength, service.HeaderContentType, service.HeaderReferer)

	return cors.New(corsConf)
}

type wildcard struct {
	prefix string
	suffix string
}

func wildcardMatchFunc(allowedOrigins []string) func(string) bool {

	var allowedWOrigins []wildcard
	allowedOriginsAll := false

	for _, origin := range allowedOrigins {
		origin = strings.ToLower(origin)
		if origin == "*" {
			allowedOriginsAll = true
			allowedWOrigins = nil
			break
		} else if i := strings.IndexByte(origin, '*'); i >= 0 {
			w := wildcard{origin[0:i], origin[i+1:]}
			allowedWOrigins = append(allowedWOrigins, w)
		}
	}

	return func(origin string) bool {
		if allowedOriginsAll {
			return true
		}

		for _, w := range allowedWOrigins {
			if w.match(origin) {
				return true
			}
		}

		return false
	}
}

func (w wildcard) match(s string) bool {
	return len(s) >= len(w.prefix+w.suffix) && strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}
