package main

import (
	"log"
	"net/http"
	"strconv"
	"io/ioutil"
	"fmt"
	"./bird"
	"./nosql"
	"encoding/json"
	"net/url"
)
const(
	unknown = iota
	add 
	delete 
	getBirdDetails
	list
)

func main() {
	log.Printf("Server started..\n")

	http.HandleFunc("/saltside", birdHandler)
	http.ListenAndServe(":8010", nil)
}
/*Methos to add bird in nosql
input
r =>request
output 
error => if fails will return error
json string of added bird
*/
func addBird(r *http.Request) ([]byte,error) {

	var newBird bird.Bird
	postData, err := ioutil.ReadAll(r.Body)
	
	if err != nil {
		log.Printf("Error: Unable to read data from request\n")
		return nil,fmt.Errorf("Unable to read data from request")
	}

	err = bird.Parse(postData,&newBird)
	if err != nil {
		log.Printf("Error: Unable to parse json\n")
		return nil,fmt.Errorf("Unable to json data")
	}

	if newBird.IsValid() == false {
		log.Printf("Error: bird mandatory param missing\n")
		return nil,fmt.Errorf("bird mandatory param missing")
	}

	newBird.SetDefault()
	value, err := json.Marshal(&newBird)
	if err == nil {
		return value,nosql.GetHandle().Add(*newBird.ID, value)
	}

	return nil,err
}
/*Methos to delete bird in nosql
id =>id of the bird to be deleted
output 
error => if fails will return error
*/
func deleteBird(id string)error {
	return nosql.GetHandle().Delete(id)
}
/*list of all birds
This is currently not supported as scan all keys is not supported by memcache
*/
func getAll()([]byte,error){

/*	birdList,err := nosql.GetHandle().GetAllKey()
	if err != nil {
		log.Printf("Error: Unable to fetch bird list\n")
		return nil,err
	}
*/
	return []byte("scanning all record is not supported by memcache\n"),nil
}
/*
Method to get details of a specific bird
input
id => id of the bird whose details is required
output
json byte slice stored at memcahe if found otherwise error
*/
func getBirdByID(id string)([]byte,error){

	birdData,err := nosql.GetHandle().Get(id)
	if err != nil{
		return nil,err
	}

	return birdData,nil
}
/*
handler for all call comming from outside
*/
func birdHandler(w http.ResponseWriter, r *http.Request) {

	var err error=nil
	var finalData []byte
	var id string
	headerResponse := 400
	
	action := getAction(r.URL)
	if action == unknown {
		log.Printf("Info: type of action is unknown")
		w.WriteHeader(400)
	}

	switch action {
	case add:
		headerResponse = 400
		finalData,err = addBird(r)
		if err != nil {
			headerResponse = 201
			finalData = nil
		}
		
	case delete:
		id,err = getBirdId(r.URL)
		headerResponse = 404
		if err == nil {
			err = deleteBird(id)
			if err == nil {
				headerResponse = 200
			}
		}
	case getBirdDetails:
		id,err = getBirdId(r.URL)
		headerResponse = 404
		if err == nil {
			finalData,err = getBirdByID(id)
			if err == nil {
				headerResponse = 200
			}
		}
			
	case list:
		finalData,err = getAll()
		headerResponse = 200
	default:
		err = fmt.Errorf("Invalid option found")
			
	}
	if err == nil {
		w.WriteHeader(headerResponse)
		if finalData != nil {
			w.Header().Add("Content-Type", "application/json")
			w.Write(finalData)
		}
	}
}
/*
Method to get bird ID from url
*/
func getBirdId(u *url.URL)(string,error){
	v := u.Query()
	id := v.Get("id")
	if len(id) > 0 {
		return id,nil
	}

	return "", fmt.Errorf("ID not found")
}
/*
Method to get type of action need to be taken
*/

func getAction(u *url.URL) int {
	
	v := u.Query()
	ops := v.Get("ops")
	if len(ops) > 0 {
		action,_ :=  strconv.Atoi(ops)
		return action
	}

	return unknown
}
