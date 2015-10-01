package main

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	baseUrlLeft := "https://query.yahooapis.com/v1/public/yql?q=select%20LastTradePriceOnly%20from%20yahoo.finance%0A.quotes%20where%20symbol%20%3D%20%22"
	baseUrlRight := "%22%0A%09%09&format=json&env=http%3A%2F%2Fdatatables.org%2Falltables.env"

	//request http api
	resp, err := http.Get(baseUrlLeft + "GOOG" + baseUrlRight)

	if err != nil {
		log.Fatal(err)
	}

	//read body
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	//check if query is successful
	if resp.StatusCode != 200 {
		log.Fatal("Query failure, possibly no network connection or illegal stock quote ")
	}

	//convert the []body into a NewJson object
	newjson, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println(err)
	}

	price, _ := newjson.Get("query").Get("results").Get("quote").Get("LastTradePriceOnly").String()
	floatPrice, err := strconv.ParseFloat(price, 64)

	fmt.Println(floatPrice)
}
