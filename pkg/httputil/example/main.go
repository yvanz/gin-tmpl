package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/yvanz/gin-tmpl/pkg/httputil"
)

type PostParams struct {
	Hello string
	World string
}

func main() {
	url := "https://www.httpbin.org/post"
	a := PostParams{
		Hello: "a",
	}

	sendBody, err := json.Marshal(&a)
	if err != nil {
		fmt.Println(err)
	}

	response, err := httputil.Post(
		url, httputil.SendBody(bytes.NewReader(sendBody)),
		httputil.SendHeaders(map[string]string{"Content-Type": "application/json", "User-Agent": "tools/1.0"}),
		httputil.SendTimeout(10*time.Second),
		httputil.SendAcceptedCodes(200, 400),
	)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	bodyByte, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("response is: %s", string(bodyByte))
}
