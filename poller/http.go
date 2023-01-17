package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// return the body as a string and as a map and the http response
func callHttp(req *http.Request) (string, map[string]interface{}, *http.Response) {
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	stringBody := string(body)

	if stringBody == "" {
		return stringBody, nil, res
	}

	// this will fail if body is nil
	var data map[string]interface{}
	err := json.Unmarshal([]byte(stringBody), &data)
	if err != nil {
		fmt.Println(err)
		fmt.Println(stringBody)
		fmt.Println(res)
		os.Exit(1)
	}
	return stringBody, data, res
}

func checkOsExit(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
		// TODO how to cause the program to exit?
	}
}