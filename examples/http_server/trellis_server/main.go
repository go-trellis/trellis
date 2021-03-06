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

package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iTrellis/trellis/cmd"
	"github.com/iTrellis/trellis/examples/components"
	"github.com/iTrellis/trellis/internal/gin_middlewares"
	"github.com/iTrellis/trellis/service"

	_ "github.com/iTrellis/trellis/server/api"
)

// curl -X 'POST' -H 'X-Api: trellis.ping' 'http://localhost:8080/v1' -H 'Authorization: aaa'

// curl -X 'GET' 'http://localhost:8080/debug/pprof/profile' -H 'Authorization: test'

// curl -i 'http://localhost:8080/' -H 'Authorization: aaa'  ## 302
// curl -i 'http://localhost:8080/text' -H 'Authorization: aaa'

func main() {
	c, err := cmd.New()
	if err != nil {
		log.Fatalln(err)
	}
	if err := c.Init(cmd.ConfigFile("config.yaml")); err != nil {
		log.Fatalln(err)
	}

	// Explicit to register component function
	cmd.DefaultCompManager.RegisterComponentFunc(
		&service.Service{Name: "component_ping", Version: "v1"},
		components.NewPing)

	gin_middlewares.RegistUseFuncs("auth", Auth())

	if err := c.BlockRun(); err != nil {
		log.Fatalln(err)
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Authorization") != "aaa" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
