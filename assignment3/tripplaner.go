package main

//@Author: Michael Dong
//@2015-10-24
//@SJSU

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	// "math/rand"
	"net/http"
	// "net/url"
	"assignment3/permutation"
	"strconv"
	"strings"
)

const (
	TOKEN        = "xxxxxxxxxx"
	ACCESS_TOKEN = "xxxxxxxxxx"
	UberBaseURL  = "https://sandbox-api.uber.com/v1/"
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

type TripRequest struct {
	Starting_from_location_ID string   `json:"starting_from_location_ID"`
	Location_IDs              []string `json:"location_IDs"`
}

type Trip struct {
	Id                            int      `json:"id" bson:"_id"`
	Request_ID                    string   `json:"request_id" bson:"request_id"`
	Status                        string   `json:"status" bson:"status"`
	Starting_from_location_ID     string   `json:"starting_from_location_ID" bson:"starting_from_location_ID"`
	Next_destionation_location_ID string   `json:"next_destionation_location_ID" bson:"next_destionation_location_ID"`
	Best_route_location_IDs       []string `json:"best_route_location_IDs" bson:"best_route_location_IDs"`
	Total_uber_costs              int      `json:"total_uber_costs" bson:"total_uber_costs"`
	Total_uber_duration           int      `json:"total_uber_duration" bson:"total_uber_duration"`
	Total_distance                float64  `json:"total_distance" bson:"total_distance"`
	Uber_wait_time_eta            int      `json:"uber_wait_time_eta" bson:"uber_wait_time_eta"`
}

type UberProduct struct {
	Product_id       string  `json:"product_id" bson:"product_id"`
	Currency_code    string  `json:"currency_code" bson:"currency_code"`
	Display_name     string  `json:"display_name" bson:"display_name"`
	Estimate         string  `json:"estimate" bson:"estimate"`
	Low_estimate     int     `json:"low_estimate" bson:"low_estimate"`
	High_estimate    int     `json:"high_estimate" bson:"high_estimate"`
	Surge_multiplier int     `json:"surge_multiplier" bson:"surge_multiplier"`
	Duration         int     `json:"duration" bson:"duration"`
	Distance         float64 `json:"distance" bson:"distance"`
}

type UberResponse struct {
	RequestID       string  `json:"request_id"`
	Status          string  `json:"status"`
	Vehicle         string  `json:"vehicle"`
	Driver          string  `json:"driver"`
	Location        string  `json:"location"`
	ETA             int     `json:"eta"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
}

type UberLocation struct {
	Lat float64
	Lng float64
}

type PutRequest struct {
	Status string `json:"status" bson:"status"`
}

type UserRequest struct {
	Product_id      string  `json:"product_id" bson:"product_id"`
	Start_latitude  float64 `json:"start_latitude" bson:"start_latitude"`
	Start_longitude float64 `json:"start_longitude" bson:"start_longitude"`
	End_latitude    float64 `json:"end_latitude" bson:"end_latitude"`
	End_longitude   float64 `json:"end_longitude" bson:"end_longitude"`
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

	if location.Address == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", "Invalid parameters")
		return
	}

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
	if location.Address == "" || location.City == "" || location.State == "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", "Invalid parameters")
		return
	}
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

func getUberProduct(startLocation UberLocation, endLocation UberLocation, uberProduct *UberProduct) {
	serviceTypeString := "estimates/price"
	parameterString := "server_token=" + TOKEN + "&" + "start_latitude=" + strconv.FormatFloat(startLocation.Lat, 'f', 6, 32) + "&" + "start_longitude=" + strconv.FormatFloat(startLocation.Lng, 'f', 6, 32) + "&" + "end_latitude=" + strconv.FormatFloat(endLocation.Lat, 'f', 6, 32) + "&" + "end_longitude=" + strconv.FormatFloat(endLocation.Lng, 'f', 6, 32)

	queryURL := UberBaseURL + serviceTypeString + "?" + parameterString

	resp, err := http.Get(queryURL)
	if err != nil {
		fmt.Println("getUberEstimate: query error, possibly parameter input error")
	}

	if resp.StatusCode != 200 {
		log.Fatalf("getUberEstimate: status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	newjson, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(body)
	// cost, err := newjson.Get("prices").GetIndex(0).Get("low_estimate").Int()
	// duration, err := newjson.Get("prices").GetIndex(0).Get("duration").Int()
	// distance, err := newjson.Get("prices").GetIndex(0).Get("distance").Float64()
	// uberEstimate.Cost = cost
	// uberEstimate.Duration = duration
	// uberEstimate.Distance = distance

	products := newjson.Get("prices")
	i := 0
	for currency, _ := products.GetIndex(i).Get("currency_code").String(); currency != "USD"; i++ {
		i++
	}
	if i != 0 {
		i--
	}

	lowest_price_product := products.GetIndex(i)
	lowest_estimate, _ := lowest_price_product.Get("low_estimate").Int()
	for ; i < len(products.MustArray()); i++ {
		currency, _ := products.GetIndex(i).Get("currency_code").String()
		current_estimate, _ := products.GetIndex(i).Get("low_estimate").Int()
		if currency == "USD" && current_estimate < lowest_estimate {
			lowest_price_product = products.GetIndex(i)
		}
	}

	productMap := lowest_price_product.Interface()

	marshaled, _ := json.Marshal(productMap)
	json.Unmarshal(marshaled, &uberProduct)

}

func getUberPrice(startId int, endId int) int {
	var start, end Location

	// search for start location
	if err := session.DB("cmpe273").C("assignment2").FindId(startId).One(&start); err != nil {
		log.Fatal("Can not get location from start ID")
	}

	// search for end location
	if err := session.DB("cmpe273").C("assignment2").FindId(endId).One(&end); err != nil {
		log.Fatal("Can not get location from start ID")

	}
	var startLocation, endLocation UberLocation
	startLocation.Lat = start.Coordinate.Lat
	startLocation.Lng = start.Coordinate.Lng
	endLocation.Lat = end.Coordinate.Lat
	endLocation.Lng = end.Coordinate.Lng

	var uberProduct UberProduct
	getUberProduct(startLocation, endLocation, &uberProduct)
	return uberProduct.Low_estimate
}

func computeBestRoute(sourceArray []string, startId string) []string {
	if len(sourceArray) == 0 {
		return []string{startId}
	}
	sourceArray = append(sourceArray)
	var sourceIntArray []int
	for i := 0; i < len(sourceArray); i++ {
		intval, _ := strconv.Atoi(sourceArray[i])
		sourceIntArray = append(sourceIntArray, intval)
	}
	permu := permutation.GeneratePermutation(sourceIntArray)

	startInt, _ := strconv.Atoi(startId)

	sourceIntArray = append(sourceIntArray, startInt)
	priceMap := computeCombination(sourceIntArray)

	bestPrice := 100000
	bestRoute := permu[0].([]int)
	for i := 0; i < len(permu); i++ {
		route := permu[i].([]int)
		routeLen := len(route)
		price := priceMap[startId+"->"+strconv.Itoa(route[0])]

		for j := 0; j < routeLen-1; j++ {
			price += priceMap[strconv.Itoa(route[j])+"->"+strconv.Itoa(route[j+1])]
		}

		price += priceMap[strconv.Itoa(route[routeLen-1])+"->"+startId]

		if price < bestPrice {
			bestPrice = price
			bestRoute = route
		}
	}

	var bestRouteStrArray []string
	for i := 0; i < len(bestRoute); i++ {
		str := strconv.Itoa(bestRoute[i])
		bestRouteStrArray = append(bestRouteStrArray, str)
	}

	return bestRouteStrArray

}

func computeCombination(sourceArray []int) map[string]int {
	if len(sourceArray) == 0 {
		return nil
	}

	result := make(map[string]int)

	for i := 0; i < len(sourceArray); i++ {
		for j := i; j < len(sourceArray); j++ {
			strFwd := strconv.Itoa(sourceArray[i]) + "->" + strconv.Itoa(sourceArray[j])
			result[strFwd] = getUberPrice(sourceArray[i], sourceArray[j])

			strBack := strconv.Itoa(sourceArray[j]) + "->" + strconv.Itoa(sourceArray[i])
			result[strBack] = getUberPrice(sourceArray[j], sourceArray[i])
		}
	}

	return result
}

// post a new trip
func postTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var tripRequest TripRequest
	json.NewDecoder(r.Body).Decode(&tripRequest)
	// fmt.Println(tripRequest)

	if tripRequest.Starting_from_location_ID == "" || len(tripRequest.Location_IDs) == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", "Invalid parameters")
		return
	}
	path := make([]string, len(tripRequest.Location_IDs)+2)
	path[0] = tripRequest.Starting_from_location_ID

	for i := 0; i < len(tripRequest.Location_IDs); i++ {
		path[i+1] = tripRequest.Location_IDs[i]
	}
	path[len(path)-1] = tripRequest.Starting_from_location_ID
	totalCost, totalDuration, totalDistance := 0, 0, 0.0

	for i := 0; i < len(path)-1; i++ {
		fmt.Println(path[i])
		startID, err := strconv.Atoi(path[i])
		endID, err := strconv.Atoi(path[i+1])

		if err != nil {
			log.Fatalf("Unable to transform id from string into int type")
		}

		var start, end Location

		if err := session.DB("cmpe273").C("assignment2").FindId(startID).One(&start); err != nil {
			log.Fatal("error retrieving location")
		}

		if err := session.DB("cmpe273").C("assignment2").FindId(endID).One(&end); err != nil {
			log.Fatal("error retrieving location")
		}

		var startLocation, endLocation UberLocation
		startLocation.Lat = start.Coordinate.Lat
		startLocation.Lng = start.Coordinate.Lng
		endLocation.Lat = end.Coordinate.Lat
		endLocation.Lng = end.Coordinate.Lng

		var product UberProduct
		getUberProduct(startLocation, endLocation, &product)

		totalCost += product.Low_estimate
		totalDistance += product.Distance
		totalDuration += product.Duration

	}

	var trip Trip

	trip.Starting_from_location_ID = tripRequest.Starting_from_location_ID
	trip.Best_route_location_IDs = computeBestRoute(tripRequest.Location_IDs, tripRequest.Starting_from_location_ID)
	trip.Next_destionation_location_ID = trip.Best_route_location_IDs[0]
	trip.Status = "planning"
	trip.Total_distance = totalDistance
	trip.Total_uber_costs = totalCost
	trip.Total_uber_duration = totalDuration

	//generate a unique autoincreasing id
	var lastId AutoIncreaseID
	session.DB("cmpe273").C("counters").Find(bson.M{}).One(&lastId)
	trip.Id = lastId.SeqToIncrease + 1

	session.DB("cmpe273").C("assignment3").Insert(trip)
	session.DB("cmpe273").C("counters").Update(bson.M{}, bson.M{"$set": bson.M{"seqToIncrease": trip.Id}})

	// Marshal provided interface into JSON structure
	tripJSON, _ := json.Marshal(trip)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", tripJSON)
}

//Get trip with id parameter
func getTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// parse id parameter
	id, _ := strconv.Atoi(p.ByName("tripId"))

	var trip Trip

	// search for trip
	if err := session.DB("cmpe273").C("assignment3").FindId(id).One(&trip); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	tripJSON, _ := json.Marshal(trip)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", tripJSON)
}

// Update a trip with location id and new information
func putTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// parse id parameter
	id, _ := strconv.Atoi(p.ByName("tripId"))

	var trip Trip

	// search for trip
	if err := session.DB("cmpe273").C("assignment3").FindId(id).One(&trip); err != nil {
		w.WriteHeader(404)
		return
	}

	var start_location, next_location string
	start_location = trip.Starting_from_location_ID
	next_location = trip.Next_destionation_location_ID

	if trip.Status == "planning" {

		putTripToUber(trip, "requesting", w)

	} else {
		if trip.Status == "requesting" {
			postTripToUber(trip, start_location, next_location, w)
		} else if trip.Status == "processing" {
			putTripToUber(trip, "accepted", w)
		} else if trip.Status == "accepted" {
			putTripToUber(trip, "arriving", w)
		} else if trip.Status == "arriving" {
			putTripToUber(trip, "in_progress", w)
		} else if trip.Status == "in_progress" {
			putTripToUber(trip, "completed", w)
		} else if trip.Status == "completed" {

			if trip.Next_destionation_location_ID == trip.Starting_from_location_ID {
				putTripToUber(trip, "finished", w)
			}
			start_location = trip.Next_destionation_location_ID
			locationArray := trip.Best_route_location_IDs

			// fmt.Println("in_progress, else2")

			for i := 0; i < len(locationArray); i++ {
				if start_location == locationArray[i] {
					i_str := strconv.Itoa(i)
					fmt.Println("in_progress, i found, i=" + i_str)
					if i != len(locationArray)-1 {
						next_location = locationArray[i+1]
						fmt.Println("in_progress, i found, i=" + i_str)
						fmt.Println("i is not last one")
					} else {
						next_location = trip.Starting_from_location_ID
						fmt.Println("i is last  one")
					}

				}

			}
			trip.Next_destionation_location_ID = next_location

			postTripToUber(trip, start_location, next_location, w)
		}

	}

}

func postTripToUber(trip Trip, start_location string, next_location string, w http.ResponseWriter) {
	serviceTypeString := "requests"
	queryString := UberBaseURL + serviceTypeString

	var start, next Location
	var startLocation, endLocation UberLocation

	startID, _ := strconv.Atoi(start_location)
	nextID, _ := strconv.Atoi(next_location)

	if err := session.DB("cmpe273").C("assignment2").FindId(startID).One(&start); err != nil {
		w.WriteHeader(404)
		return
	}

	if err := session.DB("cmpe273").C("assignment2").FindId(nextID).One(&next); err != nil {
		w.WriteHeader(404)
		return
	}

	startLocation.Lat = start.Coordinate.Lat
	startLocation.Lng = start.Coordinate.Lng
	endLocation.Lat = next.Coordinate.Lat
	endLocation.Lng = next.Coordinate.Lng

	var product UberProduct
	getUberProductFromIDString(start_location, next_location, &product, w)
	var userRequest UserRequest
	userRequest.Product_id = product.Product_id
	userRequest.Start_latitude = startLocation.Lat
	userRequest.Start_longitude = startLocation.Lng
	userRequest.End_latitude = endLocation.Lat
	userRequest.End_longitude = endLocation.Lng

	reqBody, _ := json.Marshal(userRequest)
	client := &http.Client{}

	req, err := http.NewRequest("POST", queryString, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println("post Request: query error, possibly parameter input error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var uberResponse UberResponse
	json.Unmarshal(body, &uberResponse)

	trip.Status = uberResponse.Status
	trip.Uber_wait_time_eta = uberResponse.ETA
	trip.Request_ID = uberResponse.RequestID

	session.DB("cmpe273").C("assignment3").Update(bson.M{"_id": trip.Id}, bson.M{"$set": bson.M{"status": trip.Status, "next_destionation_location_ID": trip.Next_destionation_location_ID, "uber_wait_time_eta": trip.Uber_wait_time_eta, "request_id": trip.Request_ID}})

	tripJson, _ := json.Marshal(trip)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", tripJson)
}

func putTripToUber(trip Trip, status string, w http.ResponseWriter) {

	serviceTypeString := "requests"
	queryString := UberBaseURL + serviceTypeString + "/" + trip.Request_ID
	var putRequest PutRequest
	putRequest.Status = status
	reqBody, _ := json.Marshal(putRequest)
	client := &http.Client{}

	req, err := http.NewRequest("PUT", queryString, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)

	// resp, err := client.Do(req)
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// if err != nil {
	// 	fmt.Println("put Request: query error, possibly parameter input error")
	// }

	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)

	// var uberResponse UberResponse
	// json.Unmarshal(body, &uberResponse)
	// if uberResponse.Status == "" {
	// 	trip.Status = "requesting"
	// } else {
	// 	trip.Status = uberResponse.Status
	// }

	// trip.Uber_wait_time_eta = uberResponse.ETA
	// trip.Request_ID = uberResponse.RequestID

	session.DB("cmpe273").C("assignment3").Update(bson.M{"_id": trip.Id}, bson.M{"$set": bson.M{"status": status, "next_destionation_location_ID": trip.Next_destionation_location_ID, "uber_wait_time_eta": trip.Uber_wait_time_eta, "request_id": trip.Request_ID}})

	tripJson, _ := json.Marshal(trip)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", tripJson)

}

func getUberProductFromIDString(start_location string, next_location string, uberProduct *UberProduct, w http.ResponseWriter) {
	var start, next Location
	var startLocation, endLocation UberLocation

	startID, _ := strconv.Atoi(start_location)
	nextID, _ := strconv.Atoi(next_location)

	if err := session.DB("cmpe273").C("assignment2").FindId(startID).One(&start); err != nil {
		w.WriteHeader(404)
		return
	}

	if err := session.DB("cmpe273").C("assignment2").FindId(nextID).One(&next); err != nil {
		w.WriteHeader(404)
		return
	}

	startLocation.Lat = start.Coordinate.Lat
	startLocation.Lng = start.Coordinate.Lng
	endLocation.Lat = next.Coordinate.Lat
	endLocation.Lng = next.Coordinate.Lng

	getUberProduct(startLocation, endLocation, uberProduct)
}

func main() {

	// a := []int{1, 2, 3, 4, 5}
	// var b []interface{}
	// b = permutation.GeneratePermutation(a)
	// fmt.Println(b)

	// initId = 12345

	// local database
	// session, _ = mgo.Dial("mongodb://localhost:27017")

	//remote database from mongolab
	//cmd: mongo ds039484.mongolab.com:39484/cmpe273 -u xxx -p xxx
	session, _ = mgo.Dial("mongodb://xxxxxxxx:xxxxxxx@ds039484.mongolab.com:39484/cmpe273")

	defer session.Close()

	mux := httprouter.New()
	mux.GET("/locations/:locationId", getLocation)
	mux.POST("/locations/", postLocation)
	mux.DELETE("/locations/:locationId", delLocation)
	mux.PUT("/locations/:locationId", putLocation)
	mux.POST("/trips/", postTrip)
	mux.GET("/trips/:tripId", getTrip)
	mux.PUT("/trips/:tripId/request", putTrip)
	// mux.GET("/price/", getPrice)

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
