package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"trellis.tech/trellis.v1/pkg/message"
	"trellis.tech/trellis.v1/pkg/mime"
	"trellis.tech/trellis.v1/pkg/service"
	"trellis.tech/trellis/common.v1/errcode"
	"trellis.tech/trellis/common.v1/json"
)

var (
	c  = http.Client{}
	mu = sync.Mutex{}
)

func main() {
	cc := map[string]int{}
	ch := make(chan int, 100)
	for i := 0; i < 10000; i++ {
		ch <- i
		go func(i int) {
			r := Call()
			mu.Lock()
			cc[r]++
			mu.Unlock()
			<-ch
		}(i)
	}

	time.Sleep(time.Second)
	fmt.Println(cc)
}

func Call() string {
	s := service.NewService("trellis", "componenta", "v1")
	s.Topic = "grpc"

	body := map[string]interface{}{
		"service": s,
		"name":    "haha",
		"age":     1,
	}

	bs, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1", bytes.NewBuffer(bs))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set(mime.HeaderKeyContentType, mime.ContentTypeJson)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errcode.Newf("code not 200, but %d", resp.StatusCode)
		log.Println(err)
		return err.Error()
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	response := message.Response{}
	_ = json.Unmarshal(respBody, &response)
	if response.GetPayload().Get(mime.HeaderKeyContentType) == mime.ContentTypeJson {
		r := map[string]interface{}{}
		_ = json.Unmarshal(response.GetPayload().GetBody(), &r)
		return r["message"].(string)
	}

	return string(respBody)
}
