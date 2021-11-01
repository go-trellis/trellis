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

	"trellis.tech/trellis.v0/cmd"
	"trellis.tech/trellis.v0/service"

	"trellis.tech/trellis.v0/examples/components"
	_ "trellis.tech/trellis.v0/server/http"
)

func main() {
	c, err := cmd.New()
	if err != nil {
		log.Fatalln(err)
	}
	if err := c.Init(cmd.ConfigFile("config.yaml")); err != nil {
		log.Fatalln(err)
	}

	cmd.DefaultCompManager.RegisterComponentFunc(
		&service.Service{Name: "component_ping", Version: "v1"},
		components.NewPing)

	if err := c.BlockRun(); err != nil {
		log.Fatalln(err)
	}
}
