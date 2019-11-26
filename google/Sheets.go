// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/reaction-eng/restlib/configuration"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type Sheets struct {

	//Store the connection to google.  This has been wrapped with the correct Headers
	Connection *sheets.Service
}

//Get a new interface
func NewSheets(configuation configuration.Configuration) *Sheets {
	//Create a new
	gInter := &Sheets{}

	//Open the client
	jwtConfig := &jwt.Config{
		Email:      configuation.GetStringFatal("google_auth_email"),
		PrivateKey: []byte(configuation.GetStringFatal("google_auth_key")),
		Scopes: []string{
			sheets.DriveScope,
			sheets.DriveFileScope,
			sheets.SpreadsheetsScope,
		},
		TokenURL: google.JWTTokenURL,
	}

	//Build the connection
	httpCon := jwtConfig.Client(context.Background())

	//Now build the drive service
	sheetConn, err := sheets.New(httpCon)
	gInter.Connection = sheetConn

	//Check for errors
	if err != nil {
		log.Fatalf("Unable to retrieve Sheet client: %v", err)
	}

	return gInter
}

/**
Simple method to append all of the data to the sheet
*/
func (sheetsCon *Sheets) AppendToSheet(sheetId string, sheetName string, data interface{}) error {
	//Get an empty row
	rowData, err := sheetsCon.GetEmptyDataRow(sheetId, sheetName)

	//If there was no error
	if err != nil {
		return err
	}

	//No store data in row
	rowData.StoreDataInRow(data)

	//Now append the data
	err = sheetsCon.uploadToDataSheet(rowData, sheetId, sheetName)

	return err
}

/**
Simple method to sync the headers
*/
func (sheetsCon *Sheets) SyncHeaders(sheetId string, sheetName string, headers []string) error {
	//Get an empty row
	rowData, err := sheetsCon.GetEmptyDataRow(sheetId, sheetName)

	//If there was no error
	if err != nil {
		return err
	}

	//Make a list of headers that we need to add
	headersToAdd := make([]string, 0)

	//Check each header
	for _, header := range headers {
		if !stringInSlice(header, rowData.Headers) {
			headersToAdd = append(headersToAdd, header)
		}
	}

	//If there are any
	if len(headersToAdd) > 0 {

		//Set the row to the header row
		rowData.RowNumber = 1

		//Copy the headers to data
		rowData.Values = make([]interface{}, len(rowData.Headers))

		//Copy over the data
		for c, headerName := range rowData.Headers {
			rowData.Values[c] = headerName
		}

		//Add the new headers
		for _, headerName := range headersToAdd {

			rowData.Values = append(rowData.Values, headerName)
		}

		//Now append the data
		err = sheetsCon.uploadToDataSheet(rowData, sheetId, sheetName)

	}

	return err
}

//Simple support string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/**
Append to the sheet with the name and Id split with a /
*/
func (sheetsCon *Sheets) AppendToSheetIdAndName(sheetIdAndName string, data interface{}) error {

	//Split the name with a /
	splitInfo := strings.Split(sheetIdAndName, "/")

	//If it is empty
	if len(splitInfo) == 0 {
		return errors.New("sheetid must be specified")
	} else if len(splitInfo) == 1 {
		return sheetsCon.AppendToSheet(splitInfo[0], "", data)
	} else {
		return sheetsCon.AppendToSheet(splitInfo[0], splitInfo[1], data)

	}

}

/**
Upload the row to the server
*/
func (sheetsCon *Sheets) UploadRow(sheetId string, sheetName string, row *SheetDataRow) error {
	//Check the sheet name
	sheetName, err := sheetsCon.checkSheetName(sheetId, sheetName)

	//If there was no error
	if err != nil {
		return err
	}

	//Now append the data
	err = sheetsCon.uploadToDataSheet(row, sheetId, sheetName)

	return err
}

/**
Method to get the name if not specified
*/
func (sheetsCon *Sheets) checkSheetName(sheetId string, sheetName string) (string, error) {
	//Check to see if the sheet name is specified
	if len(sheetName) > 0 {
		return sheetName, nil
	}
	//Ok, now go download the sheets
	result, err := sheetsCon.Connection.Spreadsheets.Get(sheetId).Do()

	//If there is an error
	if err != nil {
		return "", err
	}

	//Now get the result
	if len(result.Sheets) == 0 {
		return "", errors.New("no sheets available")
	}
	//Return the sheet name
	return result.Sheets[0].Properties.Title, nil

}

//Get an empty data row
func (sheetsCon *Sheets) GetEmptyDataRow(sheetId string, sheetName string) (*SheetDataRow, error) {
	//Check the sheet name
	sheetName, err := sheetsCon.checkSheetName(sheetId, sheetName)

	//If there was an error return
	if err != nil {
		return nil, err
	}

	//Get the data range
	rangeName := sheetName + "!1:1"

	//Now get the rows
	rowData, err := sheetsCon.Connection.Spreadsheets.Values.Get(sheetId, rangeName).Do()

	//If there was an error return
	if err != nil {
		return nil, err
	}

	//Make sure there are headers
	if len(rowData.Values) == 0 {
		return nil, errors.New("missing headers from sheet " + sheetName)
	}

	//Now build a new empty row
	return newEmptyDataRow(rowData.Values[0]), nil

}

//Get an empty data row
func (sheetsCon *Sheets) getMaxRows(sheetId string, sheetName string) (int, error) {

	//Get the data range
	rangeName := sheetName + ""

	//Now get the rows
	rowData, err := sheetsCon.Connection.Spreadsheets.Values.Get(sheetId, rangeName).Do()

	//If there was an error return
	if err != nil {
		return 0, err
	}

	return len(rowData.Values) + 1, nil

}

//Get an empty data row
func (sheetsCon *Sheets) GetSheetData(sheetId string, sheetName string) (*SheetData, error) {
	//Check the sheet name
	sheetName, err := sheetsCon.checkSheetName(sheetId, sheetName)

	//If there was an error return
	if err != nil {
		return nil, err
	}

	//Get the data range
	rangeName := sheetName + ""

	//Now get the rows
	rowData, err := sheetsCon.Connection.Spreadsheets.Values.Get(sheetId, rangeName).Do()

	//If there was an error return
	if err != nil {
		return nil, err
	}

	//Create a new
	sheetData, err := NewSheetData(rowData.Values)
	if err != nil {
		return nil, err
	}

	return sheetData, nil
}

//Get an empty data row
func (sheetsCon *Sheets) uploadToDataSheet(data *SheetDataRow, sheetId string, sheetName string) error {

	//If the current row is unspecified upload
	if data.RowNumber <= 0 {
		//Just append
		//Get the max number of rows
		maxRows, err := sheetsCon.getMaxRows(sheetId, sheetName)
		if err != nil {
			return err
		}

		//Get the data range
		rangeName := sheetName + "!A" + strconv.Itoa(maxRows)

		// How the input data should be interpreted.
		valueInputOption := "RAW"

		// How the input data should be inserted.
		insertDataOption := "INSERT_ROWS"

		//Build the data
		rb := &sheets.ValueRange{
			Values: make([][]interface{}, 0),
		}

		//Now store the data
		rb.Values = append(rb.Values, data.Values)

		//Now push the updated values
		_, err = sheetsCon.Connection.Spreadsheets.Values.Append(sheetId, rangeName, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Do()

		return err
	} else {
		//Replace the values
		//Get the data range
		rangeName := sheetName + "!A" + strconv.Itoa(data.RowNumber)

		// How the input data should be interpreted.
		valueInputOption := "RAW"

		//Build the data
		rb := &sheets.ValueRange{
			Values: make([][]interface{}, 0),
		}

		//Now store the data
		rb.Values = append(rb.Values, data.Values)

		//Now push the updated values
		_, err := sheetsCon.Connection.Spreadsheets.Values.Update(sheetId, rangeName, rb).ValueInputOption(valueInputOption).Do()

		return err
	}

}
