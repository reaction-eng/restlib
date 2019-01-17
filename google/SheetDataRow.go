package google

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

/**
Used to store a row of data
*/
type SheetDataRow struct {

	//Store the values in the row
	Values []interface{}

	//Store the Headers
	Headers []string

	//Store a map from the Headers to location
	headersLoc map[string]int

	//Store the row number, -1 if new
	RowNumber int
}

/**
Provide a simple function to remove white space and to lower case before comparing
*/
func SanitizeHeader(header string) string {

	//Remove duplicate white space
	space := regexp.MustCompile(`\s+`)
	header = space.ReplaceAllString(header, " ")

	//Now trim
	header = strings.TrimSpace(header)

	//Now to lower cawe
	header = strings.ToLower(header)

	return header
}

//Build an empty row based upon the Headers
func newEmptyDataRow(orgHeaders []interface{}) *SheetDataRow {
	//Get the size
	size := len(orgHeaders)

	//Now build the
	dataRow := &SheetDataRow{
		Values:     make([]interface{}, size),
		Headers:    make([]string, size),
		headersLoc: make(map[string]int, 0),
		RowNumber:  -1,
	}

	//Now store each of the Headers location for easy look up
	for loc, name := range orgHeaders {
		//Convert the name to a string
		nameString := fmt.Sprint(name)

		//Save the info
		dataRow.headersLoc[SanitizeHeader(nameString)] = loc
		dataRow.Headers[loc] = nameString
	}

	return dataRow
}

/**
Store the empty data in the row
*/
func (row *SheetDataRow) StoreDataInRow(data interface{}) error {
	//Convert to json
	jsonData, err := convertToJson(data)

	//Check for error
	if err != nil {
		return err
	}

	//Now march over map
	row.storeData(*jsonData)

	return nil

}

/**
Store the empty data in the row
*/
func (row *SheetDataRow) EraseData() {

	//Fill all of the values with an empty string
	for idx := range row.Values {
		row.Values[idx] = ""
	}

}

/**
Store the empty data in the row
*/
func (row *SheetDataRow) storeData(thisData map[string]interface{}) {
	//March over each item in map
	for key, value := range thisData {
		//Check if the value is is another map
		if asMap, isMap := value.(map[string]interface{}); isMap {
			//If it is a value map
			row.storeData(asMap)
		} else if asBool, isBool := value.(bool); isBool {
			//Now get the sanHeader
			sanHeader := SanitizeHeader(key)

			//Now store it if it is in the file
			col, inMap := row.headersLoc[sanHeader]

			//If in the map store it
			if inMap {
				if asBool {
					row.Values[col] = "x"
				} else {
					row.Values[col] = ""
				}
			}
		} else {
			//Treat it as a string and store in it the map
			valueString := fmt.Sprint(value)

			//Now get the sanHeader
			sanHeader := SanitizeHeader(key)

			//Now store it if it is in the file
			col, inMap := row.headersLoc[sanHeader]

			//If in the map store it
			if inMap {
				row.Values[col] = valueString
			}

		}

	}

}

/**
 * Looks up the value
 */
func (row *SheetDataRow) GetValue(s string) interface{} {
	col, found := row.headersLoc[s]

	if !found {
		return nil
	} else {
		if col < len(row.Values) {
			return row.Values[col]
		} else {
			return nil
		}

	}

}

/**
 * Looks up the value
 */
func (row *SheetDataRow) GetValueAsString(s string) string {
	value := row.GetValue(s)

	//If not nil return
	if value != nil {
		return fmt.Sprint(value)
	} else {
		return ""
	}

}

//Converts to a map so we can user it
func convertToJson(data interface{}) (*map[string]interface{}, error) {
	//Convert to json
	jsonBytes, err := json.Marshal(data)

	//If there was no error
	if err != nil {
		return nil, err
	}

	//Now create a result map
	jsonMap := make(map[string]interface{}, 0)

	//Now convert from the json bytes back
	err = json.Unmarshal(jsonBytes, &jsonMap)

	//Return
	return &jsonMap, err

}
