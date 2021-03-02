package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	sheets "google.golang.org/api/sheets/v4"
)

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Customers struct {
	ID     int    `json:"id"`
	Name   string `json:"Name"`
	Number string `json:"Number"`
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

	SPREADSHEETID := os.Getenv("SPREADSHEETID")

	// fmt.Println(SPREADSHEETID)
	// grab specific cells
	readRange := "B2:E"
	resp, err := srv.Spreadsheets.Values.Get(SPREADSHEETID, readRange).Do()
	// fmt.Println(resp)
	checkError(err)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// row = interface

			// r, okay := row[0].(string)

			// fmt.Println(r, okay)

			fmt.Printf("%s, %s\n", row[0], row[2])

			spreadsheet_data := Customers{
				ID:     counter,
				Name:   row[0].(string), // applies type interface of string
				Number: row[2].(string),
			}

			// writeJSON(row[0], row[5])
			customer_data = append(customer_data, spreadsheet_data)

			writeJSON(customer_data)

			counter++
		}

	}

}
func writeJSON(data []Customers) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("unable to create json file")
		return
	}

	_ = ioutil.WriteFile("customer_data.json", file, 0644)

}
