package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Response struct {
	Key   int
	Value string
}

var storeMap map[int]string

//Get location with id parameter
func getKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	var resp Response
	var resps []Response

	key_id := p.ByName("key_id")
	value := p.ByName("value")
	if key_id == "" && value == "" {
		i := 0
		for k, v := range storeMap {
			var resp Response
			resp.Key = k
			resp.Value = v
			resps[i] = resp
			i++
		}
	}

	key, err := strconv.Atoi(key_id)
	if err != nil {

	}

	var location Location

	// search for location
	if err := session.DB("cmpe273").C("assignment2").FindId(id).One(&location); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	locationJSON, _ := json.Marshal(location)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", locationJSON)
}

func putKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	id, _ := strconv.Atoi(p.ByName("locationId"))

	var location Location

	// search for location
	if err := session.DB("cmpe273").C("assignment2").FindId(id).One(&location); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	locationJSON, _ := json.Marshal(location)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", locationJSON)
}

//print usage information
func usage() {

	fmt.Println("Usage: ", os.Args[0], "portNumber")

}

func main() {
	// initId = 12345

	var port string
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments!")
		usage()
		return

	} else { //case for checking account

		portInt, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil || portInt > 65535 || portInt < 1000 {
			fmt.Println("Illegal argument! Port should be an integer between 1000 and 65535.")
			usage()
			return
		}

		portString, _ := strconv.Atoi(portInt)
		port = portString
	}

	storeMap = make([int]string)

	mux := httprouter.New()
	mux.GET("/keys/:key_id", getKey)
	mux.PUT("/keys/:key_id/:value", putKey)

	server := http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}
	server.ListenAndServe()
}
