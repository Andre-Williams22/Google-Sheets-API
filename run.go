package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// func writeJSON(data ) {
// 	file, err := json.MarshalIndent(data, "", " ")
// 	if err != nil {
// 		log.Println("unable to create json file")
// 		return
// 	}

// 	_ = ioutil.WriteFile("data.json", file, 0644)

// }

func main() {
	data, err := ioutil.ReadFile("secret.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	checkError(err)

	client := conf.Client(context.TODO())
	srv, err := sheets.New(client)
	checkError(err)

	// spreadsheetID := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	spreadsheetID := "1fpHkn3ascb5sTCuSENTM1f9-qS98CTDnhBQkT0KSd0E"

	readRange := "Sheet1!C2:40"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	// fmt.Println(resp)
	checkError(err)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			fmt.Printf("%s, %s\n", row[0], row[5])
			// writeJSON(row[0], row[5])

		}
	}
}
