package main

//@Author: Michael Dong
//@2015-10-24
//@SJSU

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
	"strconv"
	"strings"
)

type AutoIncreaseID struct {
	SeqToIncrease int `json:"seqToIncrease" bson:"seqToIncrease"`
}

type Location struct {
	Id         int    `json:"id" bson:"_id"`
	Name       string `json:"name" bson:"name"`
	Address    string `json:"address" bson:"address"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	Zip        string `json:"zip" bson:"zip"`
	Coordinate struct {
		Lat float64 `json:"lat" bson:"lat"`
		Lng float64 `json:"lng" bson:"lng"`
	} `json:"coordinate" bson:"coordinate"`
}

var session *mgo.Session

// var initId int

//Get location with id parameter
func getLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

// post a new location with user name and address
func postLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var location Location
	json.NewDecoder(r.Body).Decode(&location)
	queryString := buildQueryString(location.Address, location.City, location.State)
	getGoogleMapLocation(&location, queryString)

	//generate a unique autoincreasing id
	var lastId AutoIncreaseID
	session.DB("cmpe273").C("counters").Find(bson.M{}).One(&lastId)
	location.Id = lastId.SeqToIncrease + 1

	session.DB("cmpe273").C("assignment2").Insert(location)
	session.DB("cmpe273").C("counters").Update(bson.M{}, bson.M{"$set": bson.M{"seqToIncrease": location.Id}})

	// Marshal provided interface into JSON structure
	locationJSON, _ := json.Marshal(location)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", locationJSON)
}

// Update a location with location id and new information
func putLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// parse id parameter
	id, _ := strconv.Atoi(p.ByName("locationId"))

	var location Location

	json.NewDecoder(r.Body).Decode(&location)
	queryString := buildQueryString(location.Address, location.City, location.State)
	getGoogleMapLocation(&location, queryString)
	location.Id = id
	// fmt.Println(location)

	if err := session.DB("cmpe273").C("assignment2").Update(bson.M{"_id": location.Id}, bson.M{"$set": bson.M{"address": location.Address, "city": location.City, "state": location.State, "zip": location.Zip, "coordinate.lat": location.Coordinate.Lat, "coordinate.lng": location.Coordinate.Lng}}); err != nil {
		w.WriteHeader(404)
		return
	}
	if err := session.DB("cmpe273").C("assignment2").FindId(id).One(&location); err != nil {
		w.WriteHeader(404)
		return
	}
	// Marshal provided interface into JSON structure
	locationJSON, _ := json.Marshal(location)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", locationJSON)
}

//delete location entry from database
func delLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id, _ := strconv.Atoi(p.ByName("locationId"))

	if err := session.DB("cmpe273").C("assignment2").RemoveId(id); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)

}

// construct queryString
func buildQueryString(address string, city string, state string) string {

	var queryString string

	//process address
	addressSplitArray := strings.Split(address, " ")
	addSptArrLength := len(addressSplitArray)
	for i := 0; i < addSptArrLength; i++ {
		if i == addSptArrLength-1 {
			queryString = queryString + addressSplitArray[i] + ","
		} else {
			queryString = queryString + addressSplitArray[i] + "+"
		}
	}

	//process city
	citySplitArray := strings.Split(city, " ")
	citySptArrLength := len(citySplitArray)
	for i := 0; i < citySptArrLength; i++ {
		if i == citySptArrLength-1 {
			queryString = queryString + "+" + citySplitArray[i] + ","
		} else {
			queryString = queryString + "+" + citySplitArray[i]
		}
	}

	queryString = queryString + "+" + state
	return queryString
}

// get data from googleMapAPI
func getGoogleMapLocation(location *Location, queryString string) {
	// "http://maps.google.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA&sensor=false"
	queryURL := "http://maps.google.com/maps/api/geocode/json?address=" + queryString + "&sensor=false"
	resp, err := http.Get(queryURL)
	if err != nil {
		fmt.Println("getGoogleMapLocation: unable to get from googlemap", err)
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("getGoogleMapLocation: ioutil.ReadAll", err)
		panic(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Query failure, possibly no network connection ")
	}

	//convert the []body into a NewJson object
	newjson, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println(err)
	}

	location.Coordinate.Lat, _ = newjson.Get("results").GetIndex(0).Get("geometry").Get("location").Get("lat").Float64()
	location.Coordinate.Lng, _ = newjson.Get("results").GetIndex(0).Get("geometry").Get("location").Get("lng").Float64()

}

func main() {
	// initId = 12345

	//local database
	// session, _ = mgo.Dial("mongodb://localhost:27017")

	//remote database from mongolab
	//cmd: mongo ds039484.mongolab.com:39484/cmpe273 -u xxx -p xxx
	session, _ = mgo.Dial("mongodb://michael:cmpe273@ds039484.mongolab.com:39484/cmpe273")

	defer session.Close()

	mux := httprouter.New()
	mux.GET("/locations/:locationId", getLocation)
	mux.POST("/locations/", postLocation)
	mux.DELETE("/locations/:locationId", delLocation)
	mux.PUT("/locations/:locationId", putLocation)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
