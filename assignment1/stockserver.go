// Copyright 2015 Michael Dong@SJSU. All rights reserved.

package main

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/rpc"

	"strconv"
	"strings"
)

type StockAccounts struct {
	stockPortfolio map[int](*Portfolio)
}

type Portfolio struct {
	stocks           map[string](*Share)
	uninvestedAmount float32
}

type Share struct {
	boughtPrice float32
	shareNum    int
}

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

//declare a global stock account variable
var st StockAccounts

//declare a global tradeId variable, will increase by 1 per each buy
var tradeId int

//buy shares, store data into stock account and return summary
func (st *StockAccounts) Buy(rq *BuyRequest, rsp *BuyResponse) error {

	//increase tradeId by 1 per each Buy, so they will always be unique
	tradeId++
	rsp.TradeNum = tradeId

	//if account not set up, set up 1 first
	if st.stockPortfolio == nil {

		st.stockPortfolio = make(map[int](*Portfolio))

		st.stockPortfolio[tradeId] = new(Portfolio)
		st.stockPortfolio[tradeId].stocks = make(map[string]*Share)

	}

	//parse the buy arguments
	symbolAndPercentages := strings.Split(rq.StockSymbolAndPercentage, ",")
	newbudget := float32(rq.Budget)
	var spent float32

	for _, stk := range symbolAndPercentages {

		//parse how many shares and their separate budget
		splited := strings.Split(stk, ":")
		if len(splited) < 2 {
			return errors.New("Wrong trade argument! ")
		}

		stkQuote := splited[0]
		percentage := splited[1]
		strPercentage := strings.TrimSuffix(percentage, "%")
		floatPercentage64, _ := strconv.ParseFloat(strPercentage, 32)
		floatPercentage := float32(floatPercentage64 / 100.00)
		currentPrice := checkQuote(stkQuote)

		if currentPrice == 0 {
			return errors.New(stkQuote + " : no such stock name!")
		}

		shares := int(math.Floor(float64(newbudget * floatPercentage / currentPrice)))
		sharesFloat := float32(shares)
		spent += sharesFloat * currentPrice

		// if it's a new tradeId, set up a portfolio first
		if _, ok := st.stockPortfolio[tradeId]; !ok {

			newPortfolio := new(Portfolio)
			newPortfolio.stocks = make(map[string]*Share)
			st.stockPortfolio[tradeId] = newPortfolio
		}
		if _, ok := st.stockPortfolio[tradeId].stocks[stkQuote]; !ok {

			newShare := new(Share)
			newShare.boughtPrice = currentPrice
			newShare.shareNum = shares
			st.stockPortfolio[tradeId].stocks[stkQuote] = newShare
		} else {

			total := float32(sharesFloat*currentPrice) + float32(st.stockPortfolio[tradeId].stocks[stkQuote].shareNum)*st.stockPortfolio[tradeId].stocks[stkQuote].boughtPrice
			st.stockPortfolio[tradeId].stocks[stkQuote].boughtPrice = total / float32(shares+st.stockPortfolio[tradeId].stocks[stkQuote].shareNum)
			st.stockPortfolio[tradeId].stocks[stkQuote].shareNum += shares
		}

		stockBought := stkQuote + ":" + strconv.Itoa(shares) + ":$" + strconv.FormatFloat(float64(currentPrice), 'f', 2, 32)

		rsp.Stocks = append(rsp.Stocks, stockBought)
	}

	//calculate uninvested amount
	leftOver := newbudget - spent
	rsp.UninvestedAmount = leftOver
	st.stockPortfolio[tradeId].uninvestedAmount += leftOver

	return nil
}

//check account with trade number
func (st *StockAccounts) CheckAccount(args string, checkResp *CheckResponse) error {

	if st.stockPortfolio == nil {
		return errors.New("No account set up yet.")
	}

	//parse argument into a tradeId
	tradeNum64, err := strconv.ParseInt(args, 10, 64)

	if err != nil {
		return errors.New("Illegal Trade ID. ")
	}
	tradeNum := int(tradeNum64)

	if pocket, ok := st.stockPortfolio[tradeNum]; ok {

		var currentMarketVal float32
		for stockquote, sh := range pocket.stocks {
			//obtain current price
			currentPrice := checkQuote(stockquote)

			//obtain price when bought,and compare with current price to determine up or down
			var str string
			if sh.boughtPrice < currentPrice {
				str = "+$" + strconv.FormatFloat(float64(currentPrice), 'f', 2, 32)
			} else if sh.boughtPrice > currentPrice {
				str = "-$" + strconv.FormatFloat(float64(currentPrice), 'f', 2, 32)
			} else {
				str = "$" + strconv.FormatFloat(float64(currentPrice), 'f', 2, 32)
			}

			//setup object to response back
			entry := stockquote + ":" + strconv.Itoa(sh.shareNum) + ":" + str

			checkResp.Stocks = append(checkResp.Stocks, entry)

			currentMarketVal += float32(sh.shareNum) * currentPrice
		}

		//calculated uninvested amount
		checkResp.UninvestedAmount = pocket.uninvestedAmount

		//calculate total market value of holding shares
		checkResp.TotalMarketValue = currentMarketVal
	} else {
		return errors.New("No such trade ID. ")
	}

	return nil
}

func main() {

	//initialize the stock account
	var st = *(new(StockAccounts))

	//initialize a tradeId with random number
	tradeId = rand.Intn(10000) + 1

	//register the stock account data and start server with HTTP protocol
	rpc.Register(&st)
	rpc.HandleHTTP()

	//start listening
	err := http.ListenAndServe(":1234", nil) //nil, no need for handler

	checkError(err)

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func checkQuote(stockName string) float32 {
	//yahoo api, to simplify, query only one stock each time
	baseUrlLeft := "https://query.yahooapis.com/v1/public/yql?q=select%20LastTradePriceOnly%20from%20yahoo.finance%0A.quotes%20where%20symbol%20%3D%20%22"
	baseUrlRight := "%22%0A%09%09&format=json&env=http%3A%2F%2Fdatatables.org%2Falltables.env"

	//request http api
	resp, err := http.Get(baseUrlLeft + stockName + baseUrlRight)

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

	//obtain the LastTradePriceOnly, which is considered current stock price
	price, _ := newjson.Get("query").Get("results").Get("quote").Get("LastTradePriceOnly").String()
	floatPrice, err := strconv.ParseFloat(price, 32)

	// fmt.Println(floatPrice)
	return float32(floatPrice)
}
