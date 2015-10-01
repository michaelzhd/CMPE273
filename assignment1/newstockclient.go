// Copyright 2015 Michael Dong@SJSU. All rights reserved.
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
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

	//connect to server
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	checkError(err)

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

		chkResp := new(CheckResponse)

		err = client.Call("StockAccounts.CheckAccount", os.Args[1], &chkResp)

		checkError(err)

		fmt.Print("stocks: ")
		fmt.Println(chkResp.Stocks)

		fmt.Print("uninvested amount: ")
		fmt.Println(chkResp.UninvestedAmount)

		fmt.Print("total market value: ")
		fmt.Println(chkResp.TotalMarketValue)

	} else if len(os.Args) == 3 { //case for buy new shares
		budget, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			fmt.Println("Wrong budget argument.")
			usage()
			return
		}
		buyReq := new(BuyRequest)
		buyReq.StockSymbolAndPercentage = os.Args[1]
		buyReq.Budget = float32(budget)

		buyResp := new(BuyResponse)

		err = client.Call("StockAccounts.Buy", &buyReq, &buyResp)
		checkError(err)

		fmt.Print("trade Id: ")
		fmt.Println(buyResp.TradeNum)

		fmt.Print("stocks: ")
		fmt.Println(buyResp.Stocks)

		fmt.Print("uninvested amount: ")
		fmt.Println(buyResp.UninvestedAmount)

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
