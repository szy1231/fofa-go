// Package examples shows how to use fofa sdk
package examples

import (
	"fmt"
	"github.com/xiaoyu-0814/fofa-go/fofa"
	"os"
)

// FofaExample fofa sdk functons included
func FofaExample() {
	email := os.Getenv("FOFA_EMAIL")
	key := os.Getenv("FOFA_KEY")

	clt := fofa.NewFofaClient([]byte(email), []byte(key))
	if clt == nil {
		fmt.Printf("create fofa client\n")
		return
	}
	//QueryAsJSON
	ret, err := clt.QueryAsJSON(1, 100,"false", []byte(`body="小米"`))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	fmt.Printf("%s\n", ret)
	//QueryAsObject
	data, err := clt.QueryAsObject(1, 1000,"false", []byte(`domain="163.com"`), []byte("ip,host,title"))
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return
	}
	fmt.Printf("count: %d\n", len(data.Results))
	fmt.Printf("\n%s\n", data.String())
}
