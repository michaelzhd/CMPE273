package main

//@Author: Michael Dong
//@2015-11-25
//@SJSU
import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Node struct {
	NodeAddress string `json:"nodeAddress"`
}

//declare the global map to stop node address and keys
var hashMapCircle map[string]string

//declare the global nodes slice
var nodes []string

//use md5 as hash code generating function, chew up a string and return a string comprised
//of hexidecimal form of the md5 result
func hashCode(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))

}

//get all keys if the value is empty string
func getAllKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var resps []Response
	for _, nodeString := range nodes {
		node := hashMapCircle[nodeString]

		// fmt.Println(node + "/keys/")
		resp, err := http.Get(node + "/keys/")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("getAllKey: ", err)
			panic(err)
		}

		var response []Response
		json.Unmarshal(body, &response)

		//merge the results from every node into one json dictionary
		for _, res := range response {
			if res.Key != "" {
				resps = append(resps, res)
			}
		}
	}
	// Marshal provided interface into JSON structure
	respJSON, _ := json.Marshal(resps)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", respJSON)
}

//Get key with key_id parameter
func getKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	// var resp Response

	key_id := p.ByName("key_id")
	hashedKey := hashCode(key_id)
	nodeAddress := findNode(hashedKey)

	resp, err := http.Get(nodeAddress + "/keys/" + key_id)
	if err != nil {
		fmt.Println("Unable to fetch key, please check node : " + nodeAddress)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("getKey: ioutil.ReadAll", err)
		}

		var response Response
		json.Unmarshal(body, &response)
		respJSON, _ := json.Marshal(response)

		// Write content-type, status code, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", respJSON)
	} else {
		w.WriteHeader(404)
	}
}

func putKey(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	key_id := p.ByName("key_id")
	value := p.ByName("value")

	if key_id == "" {
		w.WriteHeader(400)
	}

	hashedKey := hashCode(key_id)
	nodeAddress := findNode(hashedKey)

	client := &http.Client{}
	queryString := nodeAddress + "/keys/" + key_id + "/" + value
	req, err := http.NewRequest("PUT", queryString, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Do(req)

	var response Response
	response.Key = key_id
	response.Value = value
	responseJSON, _ := json.Marshal(response)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", responseJSON)

}

func postNodes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var node Node
	json.NewDecoder(r.Body).Decode(&node)
	validAddress := regexp.MustCompile(`^(http://)(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):([1-6][0-5][0-5][0-3][0-4]|[1-9][0-9]{3})$`)
	if validAddress.MatchString(node.NodeAddress) {
		addNode(node)
		w.WriteHeader(201)
	} else {
		w.WriteHeader(400)
	}
}

func deleteNodes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var node Node
	json.NewDecoder(r.Body).Decode(&node)
	validAddress := regexp.MustCompile(`^(http://)(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]):([1-6][0-5][0-5][0-3][0-4]|[1-9][0-9]{3})$`)
	if validAddress.MatchString(node.NodeAddress) {
		removeNode(node)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

}

func addNode(node Node) {
	hashedAddess := hashCode(node.NodeAddress)
	hashMapCircle[hashedAddess] = node.NodeAddress
	nodes = append(nodes, hashedAddess)
	sort.Strings(nodes)

}

func removeNode(node Node) {
	hashedAddress := hashCode(node.NodeAddress)
	if _, ok := hashMapCircle[hashedAddress]; ok {
		delete(hashMapCircle, hashedAddress)
		for i := 0; i < len(nodes); i++ {
			if nodes[i] == hashedAddress {

				// nodes = append(nodes[:i], nodes[i+1:])
				nodes = nodes[:i+copy(nodes[i:], nodes[i+1:])]

			}
		}
	}
}

func findNode(input string) string {
	hashedAddress := hashCode(input)
	for i := 0; i < len(nodes)-1; i++ {
		if hashedAddress > nodes[i] && hashedAddress < nodes[i+1] {
			return hashMapCircle[nodes[i+1]]
		}
	}

	// if hashedAddress > nodes[len(nodes)-1] || hashedAddress < nodes[0]
	return hashMapCircle[nodes[0]]

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

	} else { //case for checking account

		portInt, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil || portInt > 65535 || portInt < 1000 {
			fmt.Println("Illegal argument! Port should be an integer between 1000 and 65535.")
			usage()
			return
		}

		port = strconv.Itoa(int(portInt))
	}

	//the global map to stop node address, the key is the hashcode of the address
	hashMapCircle = make(map[string]string)

	//initialize 3 nodes, just for testing
	addNode(Node{"http://127.0.0.1:3000"})
	addNode(Node{"http://127.0.0.1:3001"})
	addNode(Node{"http://127.0.0.1:3002"})

	mux := httprouter.New()
	mux.GET("/keys/", getAllKey)
	mux.GET("/keys/:key_id/", getKey)
	mux.PUT("/keys/:key_id/:value", putKey)
	mux.POST("/nodes/", postNodes)
	mux.DELETE("/nodes/", deleteNodes)

	server := http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: mux,
	}
	server.ListenAndServe()
}
