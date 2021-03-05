package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	sheets "google.golang.org/api/sheets/v4"
)

type Customers struct {
	ID     int    `json:"id"`
	Name   string `json:"Name"`
	Number string `json:"Number"`
}

type CustomData struct {
	CustomData []Customer `json:"customer_data"`
}

type Customer struct {
	ID     string `json:"id"`
	Name   string `json:"Name"`
	Number string `json:"Number"`
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// sendMessage sends a text to specific users
func sendMessage(to, from, message string) {

	// Load Twilio API Info

	sid := os.Getenv("SID")
	token := os.Getenv("TOKEN")

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + sid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	twilio := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(sid, token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, _ := twilio.Do(req)
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["body"])

		}
	} else {
		fmt.Println(response.Status)
	}

	// return message

}

// numPeople counts number of people in our file
func numPeople(counter int) (total int) {

	total = counter + 1

	return total
}

// writeJSON writes data to a Json file
func writeJSON(data []Customers) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("unable to create json file")
		return
	}

	_ = ioutil.WriteFile("customer_data.json", file, 0644)

}

func main() {
	// load in environment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	customer_data := make([]Customers, 0)
	// initalize counter for customer ID's
	var counter int

	// Google API Credentials
	data, err := ioutil.ReadFile("secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	checkError(err)

	client := conf.Client(context.TODO())
	srv, err := sheets.New(client)
	checkError(err)

	SPREADSHEETID := os.Getenv("TEST_SHEET")

	// fmt.Println(SPREADSHEETID)
	// grab specific cells
	readRange := "B2:E"
	resp, err := srv.Spreadsheets.Values.Get(SPREADSHEETID, readRange).Do()
	// fmt.Println(resp)
	checkError(err)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Number")
		for _, row := range resp.Values {

			fmt.Printf("%s, %s\n", row[0], row[2])

			spreadsheet_data := Customers{
				ID:     counter,
				Name:   row[0].(string), // applies type interface of string
				Number: row[2].(string),
			}

			// add data to list
			customer_data = append(customer_data, spreadsheet_data)

			// convert list to Json
			writeJSON(customer_data)

			// increment counter
			counter++
		}

	}

	// calculates total num people in our db
	numPeople(counter)

	// Read's in Customer Data
	file, err := ioutil.ReadFile("customer_data.json")

	if err != nil {
		fmt.Println(err.Error())
	}
	// create users array
	var clients []Customer

	// loads in data from json file
	err2 := json.Unmarshal(file, &clients)
	if err2 != nil {
		fmt.Println("Error Json Marshalling")
	}

	// loops through json and prints out data
	for _, x := range clients {
		fmt.Printf("Name: %s \n", x.Name)
		fmt.Printf("Number: %s \n", x.Number)
		// sends message
		sendMessage(x.Number, os.Getenv("TWILIO_NUMBER"), fmt.Sprintf("Hello %s", x.Name+"thanks for using Hancock Appliance Repair, please fill out our survey: https://www.surveymonkey.com/r/GZB6CRY"))

	}

}
