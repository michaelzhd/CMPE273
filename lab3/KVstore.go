package main

//@Author: Michael Dong
//@2015-11-25
//@SJSU

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//declare the global map to stop keys and values
var storeMap map[string]string

func getAllKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var resps []Response

	for k, v := range storeMap {
		var resp Response
		resp.Key = k
		resp.Value = v
		resps = append(resps, resp)
	}

	// Marshal provided interface into JSON structure
	respJSON, _ := json.Marshal(resps)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", respJSON)

}

//Get key with id parameter
func getKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	var resp Response

	key_id := p.ByName("key_id")

	if val, ok := storeMap[key_id]; ok {
		resp.Key = key_id
		resp.Value = val
		// Marshal provided interface into JSON structure
		respJSON, _ := json.Marshal(resp)

		// Write content-type, status code, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", respJSON)
	} else {
		w.WriteHeader(404)
	}
}

//put method of RESTful services, modify or add the value of a key
func putKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	key_id := p.ByName("key_id")
	value := p.ByName("value")
	if key_id == "" {
		w.WriteHeader(400)
	}
	storeMap[key_id] = value

	var response Response
	response.Key = key_id
	response.Value = value
	responseJSON, _ := json.Marshal(response)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", responseJSON)

}

//put method of RESTful services, delete the key if value is set to empty
func delKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	key_id := p.ByName("key_id")
	if key_id == "" {
		w.WriteHeader(400)
	}

	delete(storeMap, key_id)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
}

//print usage information
func usage() {

	fmt.Println("Usage:\n " + os.Args[0] + ":portNumber")

}

func main() {
	// initId = 12345

	var port string
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments!")
		usage()
		return

	} else {

		portInt, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil || portInt > 65535 || portInt < 1000 {
			fmt.Println("Illegal argument! Port should be an integer between 1000 and 65535.")
			usage()
			return
		}

		port = strconv.Itoa(int(portInt))
	}

	//initialize the map to store key and value data
	storeMap = make(map[string]string)

	mux := httprouter.New()
	mux.GET("/keys/", getAllKey)
	mux.GET("/keys/:key_id/", getKey)
	mux.PUT("/keys/:key_id/:value", putKey)
	mux.PUT("/keys/:key_id/", delKey)

	server := http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}
	server.ListenAndServe()
}
