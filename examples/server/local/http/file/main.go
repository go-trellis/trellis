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

package main

import (
	"fmt"

	"trellis.tech/trellis.v1/pkg/server"

	routing "github.com/go-trellis/fasthttp-routing"
	"trellis.tech/trellis.v1/cmd"
	_ "trellis.tech/trellis.v1/examples/components"
)

var (
	use1 routing.Handler = func(*routing.Context) error {
		fmt.Println("I am an use handler")

		return nil
	}
)

func init() {
	server.RegisterUseFunc("use1", use1)
}

func main() {
	cmd.Run()
}
