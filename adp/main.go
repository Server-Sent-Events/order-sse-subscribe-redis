package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type MeuTeste struct {
	merchant    string
	numberLogic string
}

func main() {

	letters := []MeuTeste{&MeuTeste{merchant: "0010489937500001", numberLogic: "78401155-3"}}
	a := &MeuTeste{merchant: "", numberLogic: ""}

	url := "http://xxxxxxxxxx/api/v2/installation-groups/5e4b1e41-4906-4873-8c9a-5fcabe15dbdb/installation-rules"

	payload := strings.NewReader("{\n  \"merchant_id\": \"0020000230310001\",\n  \"number\": \"38110024-7\",\n  \"installation_group_uuid\" : \"035c2cf9-c689-461b-92a9-776e6c2235cb\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("merchant-id", "0020000230310001")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("postman-token", "f258dcaf-243b-882d-91da-81618e00981b")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
