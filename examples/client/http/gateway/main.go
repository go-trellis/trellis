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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
)

var (
	mu  = sync.Mutex{}
	num int
)

type Response struct {
	Code    int64           `json:"code"`
	ErrMsg  string          `json:"err_msg"`
	Payload message.Payload `json:"payload"`
}

func main() {
	flag.IntVar(&num, "num", 1, "request count")
	flag.Parse()

	cc := map[string]int{}
	ch := make(chan int, 100)

	hc := &http.Client{}
	for i := 0; i < num; i++ {
		ch <- i
		go func(i int) {
			r := call(hc)
			mu.Lock()
			cc[r]++
			mu.Unlock()
			<-ch
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println(cc)
}

func call(hc *http.Client) string {
	reader := strings.NewReader(`{"name": "peter", "age": 1}`)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/v1", reader)
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Content-Type", "application/json")
	hResp, err := hc.Do(req)
	if err != nil {
		return err.Error()
	}

	defer func() {
		io.Copy(io.Discard, hResp.Body)
		hResp.Body.Close()
	}()

	if hResp.StatusCode != 200 {
		fmt.Printf("status not ok, but %d\n", hResp.StatusCode)
	}

	body, err := io.ReadAll(hResp.Body)
	if err != nil {
		return err.Error()
	}

	resp := Response{}

	if err = json.Unmarshal(body, &resp); err != nil {
		return err.Error()
	}

	if resp.Code != 0 {
		return fmt.Sprintf("%d, %s", resp.Code, resp.ErrMsg)
	}
	ct := resp.Payload.Get(mime.HeaderKeyContentType)
	if ct == mime.ContentTypeJson {
		r := map[string]interface{}{}
		_ = json.Unmarshal(resp.Payload.GetBody(), &r)
		return r["message"].(string)
	}
	return fmt.Sprintf("content-type err: %s", ct)
}
