// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"errors"
	"fmt"
	"strings"
)

/**
Used to store a row of data
*/
type SheetData struct {

	//Store the values in the row
	Values [][]interface{} //[row][col]

	//Store the Headers
	Headers []string

	//Store a map from the Headers to location
	headersLoc map[string]int

	//Store a list of original row numbers
	RowNumb []int
}

//Create a new sheet data
func NewSheetData(values [][]interface{}) (*SheetData, error) {

	//We need at least a header row
	if len(values) == 0 {
		return nil, errors.New("header row missing from sheet")
	}

	//Get the size
	headerSize := len(values[0])

	//Build the needed header info
	data := &SheetData{
		Headers:    make([]string, headerSize),
		headersLoc: make(map[string]int, 0),
		RowNumb:    make([]int, len(values)-1), //No need for the header
		Values:     values[1:],                 //Remove the header row from the data
	}

	//Now store each of the Headers location for easy look up
	for loc, name := range values[0] {
		//Convert the name to a string
		nameString := fmt.Sprint(name)

		//Save the info
		data.headersLoc[SanitizeHeader(nameString)] = loc
		data.Headers[loc] = nameString
	}

	//Now just set the values for the row location
	for i := 0; i < len(data.RowNumb); i++ {
		data.RowNumb[i] = i + 2 //C style index plus the missing header row
	}

	return data, nil
}

//March over each header and look for value
func (sheet *SheetData) findHeaderContaining(header string) int {
	for ind, value := range sheet.Headers {
		if strings.Contains(strings.ToLower(value), header) {
			return ind
		}
	}
	return -1

}

//Create a new sheet data
func (sheet *SheetData) FilterSheet(header string, value string) *SheetData {
	//See if we have the header
	headerCol := sheet.findHeaderContaining(header)

	//If it is there, return nil
	if headerCol < 0 {
		return nil
	}

	//Now create a new data sheet
	newSheet := &SheetData{
		Headers:    sheet.Headers,    //Headers are the same
		headersLoc: sheet.headersLoc, //Headers are the same
		RowNumb:    make([]int, 0),   //No need for the header
		Values:     make([][]interface{}, 0),
	}

	//Clean up the value
	value = strings.TrimSpace(value)

	//Now check to see if each row has the data
	for r, rowData := range sheet.Values {
		//Make sure we have enough data to check
		if len(rowData) > headerCol {
			//If they are equal
			if strings.EqualFold(value, strings.TrimSpace(fmt.Sprint(rowData[headerCol]))) {
				//Add the data
				newSheet.Values = append(newSheet.Values, rowData)
				newSheet.RowNumb = append(newSheet.RowNumb, sheet.RowNumb[r])
			}

		}

	}

	//Return the new sheet
	return newSheet
}

//Create a new sheet data
func (sheet *SheetData) GetColumn(col int) []interface{} {
	//We need to transpose the data
	colData := make([]interface{}, 0)

	//March over each row
	for r := range sheet.Values {
		//If it has the column add it
		if len(sheet.Values[r]) > col+1 {
			colData = append(colData, sheet.Values[r][col])
		}

	}

	return colData
}

//Create a new sheet data
func (sheet *SheetData) GetRow(row int) *SheetDataRow {
	//Look up the row number
	indexNumber := -1

	//Now search over the rows for the row index,
	for index, rowTest := range sheet.RowNumb {
		if rowTest == row {
			indexNumber = index
		}
	}

	//If it avaiable return
	if indexNumber < 0 {
		return nil
	}

	//Extract the row info
	dataRow := &SheetDataRow{
		Values:     sheet.Values[indexNumber],
		Headers:    sheet.Headers,
		headersLoc: sheet.headersLoc,
		RowNumber:  row,
	}

	//If we have fewer values then header size
	if len(dataRow.Values) < len(dataRow.Headers) {
		//Make a new array
		array := make([]interface{}, len(dataRow.Headers)-len(dataRow.Values))

		dataRow.Values = append(dataRow.Values, array...)
	}
	//Return the new sheet
	return dataRow
}

//Create a new sheet data
func (sheet *SheetData) GetEmptyDataRow() *SheetDataRow {

	//Get the size. number of headers
	size := len(sheet.Headers)

	//Extract the row info
	dataRow := &SheetDataRow{
		Values:     make([]interface{}, size),
		Headers:    sheet.Headers,
		headersLoc: sheet.headersLoc,
		RowNumber:  -1,
	}

	//Return the new sheet
	return dataRow
}

//Create a new sheet data
func (sheet *SheetData) PrintToScreen() {
	//March over each header
	fmt.Print("row,")
	for _, header := range sheet.Headers {
		fmt.Print(header)
		fmt.Print(",")
	}
	fmt.Println()

	//Now print each row
	for r, rowData := range sheet.Values {
		//Now print the row number
		fmt.Print(sheet.RowNumb[r])
		fmt.Print(",")

		//Now each data
		for _, data := range rowData {
			fmt.Print(data)
			fmt.Print(",")
		}

		fmt.Println()

	}

}

//Create a new sheet data
func (sheet *SheetData) NumberRows() int {
	return len(sheet.Values)

}

/**
Count the number of entries for this column
*/
func (sheet *SheetData) CountEntries(index int) int {
	//Start with a count
	count := 0

	//March over each row
	for _, rowData := range sheet.Values {
		//If the row as data in the index
		if index < len(rowData) {
			//If there is data
			if rowData[index] != nil && len(fmt.Sprint(rowData[index])) > 0 {
				count++
			}
		}

	}

	return count
}
