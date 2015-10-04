// Copyright 2015 Michael Dong@SJSU. All rights reserved.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "net/rpc"
	"github.com/bitly/go-simplejson"
	"os"
	"strconv"
	"strings"
)

type BuyRequest struct {
	StockSymbolAndPercentage string
	Budget                   float32
}

type BuyResponse struct {
	TradeNum         int
	Stocks           []string
	UninvestedAmount float32
}

type CheckResponse struct {
	Stocks           []string
	UninvestedAmount float32
	TotalMarketValue float32
}

func main() {

	// //connect to server
	// client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	// checkError(err)

	//parse arguments
	if len(os.Args) > 4 || len(os.Args) < 2 {
		fmt.Println("Wrong number of arguments!")
		usage()
		return

	} else if len(os.Args) == 2 { //case for checking account

		_, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			fmt.Println("Illegal argument!")
			usage()
			return
		}

		// chkResp := new(CheckResponse)

		data, err := json.Marshal(map[string]interface{}{
			"method": "StockAccounts.Check",
			"id":     1,
			"params": []map[string]interface{}{map[string]interface{}{"TradeId": os.Args[1]}},
		})

		if err != nil {
			log.Fatalf("Marshal : %v", err)
		}

		resp, err := http.Post("http://127.0.0.1:1234/rpc", "application/json", strings.NewReader(string(data)))

		if err != nil {
			log.Fatalf("Post: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalf("ReadAll: %v", err)
		}

		newjson, err := simplejson.NewJson(body)

		checkError(err)

		fmt.Print("stocks: ")
		stocks := newjson.Get("result").Get("Stocks")
		fmt.Println(*stocks)

		fmt.Print("uninvested amount: ")
		uninvestedAmount, _ := newjson.Get("result").Get("UninvestedAmount").Float64()
		fmt.Print("$")
		fmt.Println(uninvestedAmount)

		fmt.Print("total market value: ")
		totalMarketValue, _ := newjson.Get("result").Get("TotalMarketValue").Float64()
		fmt.Print("$")
		fmt.Println(totalMarketValue)

	} else if len(os.Args) == 3 { //case for buy new shares
		budget, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			fmt.Println("Wrong budget argument.")
			usage()
			return
		}

		data, err := json.Marshal(map[string]interface{}{
			"method": "StockAccounts.Buy",
			"id":     2,
			"params": []map[string]interface{}{map[string]interface{}{"StockSymbolAndPercentage": os.Args[1], "Budget": float32(budget)}},
		})

		if err != nil {
			log.Fatalf("Marshal : %v", err)
		}

		resp, err := http.Post("http://127.0.0.1:1234/rpc", "application/json", strings.NewReader(string(data)))

		if err != nil {
			log.Fatalf("Post: %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalf("ReadAll: %v", err)
		}

		newjson, err := simplejson.NewJson(body)

		checkError(err)

		fmt.Print("Trade Num: ")
		tradenum, _ := newjson.Get("result").Get("TradeNum").Int()
		fmt.Println(tradenum)

		fmt.Print("stocks: ")
		stocks := newjson.Get("result").Get("Stocks")
		fmt.Println(*stocks)

		fmt.Print("uninvested amount: ")
		uninvestedAmount, _ := newjson.Get("result").Get("UninvestedAmount").Float64()
		fmt.Print("$")
		fmt.Println(uninvestedAmount)

	} else {
		fmt.Println("Unknown error.")
		usage()
		return
	}

}

//check all kinds of error
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		log.Fatal("error: ", err)
		os.Exit(2)
	}

}

//print usage information
func usage() {

	fmt.Println("Usage: ", os.Args[0], "tradeId")
	fmt.Println("or")
	fmt.Println(os.Args[0], "“GOOG:50%,YHOO:50%” 10000(your budget)")
}
