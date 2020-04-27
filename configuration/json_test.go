// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package configuration

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJson_Get(t *testing.T) {
	testCases := []struct {
		configString   string
		key            string
		expectedResult interface{}
	}{
		{`{"testKey":true}`, "testKey", true},
		{`{"testKey":false}`, "testKey", false},
		{`{"test Key":"blue string"}`, "test Key", "blue string"},
		{`{"testKey":"blue string"}`, "testKeyNotFound", nil},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, err := NewJson(testCase.configString)
		assert.Nil(t, err)

		// act
		result := jsonConfig.Get(testCase.key)

		// assert
		assert.Equal(t, testCase.expectedResult, result)
	}

	// assert that the env takes over
	for _, testCase := range testCases {
		// arrange
		jsonConfig, err := NewJson(testCase.configString)
		assert.Nil(t, err)
		os.Setenv(testCase.key, "test value 123")

		// act
		result := jsonConfig.Get(testCase.key)

		// assert
		assert.Equal(t, "test value 123", result)

		os.Unsetenv(testCase.key)
	}
}

func TestJson_GetFatal(t *testing.T) {
	testCases := []struct {
		configString   string
		key            string
		expectedResult interface{}
		fatalMessage   string
	}{
		{`{"testKey":true}`, "testKey", true, ""},
		{`{"testKey":false}`, "testKey", false, ""},
		{`{"test Key":"blue string"}`, "test Key", "blue string", ""},
		{`{"testKey":"blue string"}`, "testKeyNotFound", nil, "Cannot not find configuration for testKeyNotFound"},
		{`{}`, "testKeyNotFound", nil, "Cannot not find configuration for testKeyNotFound"},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, err := NewJson(testCase.configString)
		assert.Nil(t, err)

		var fatalMessage string
		jsonConfig.fatal = func(i interface{}) {
			fatalMessage = fmt.Sprint(i)
		}

		// act
		result := jsonConfig.GetFatal(testCase.key)

		// assert
		assert.Equal(t, testCase.expectedResult, result)
		assert.Equal(t, testCase.fatalMessage, fatalMessage)
	}
}

func TestJson_GetString(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     string
	}{
		{`{"testKey":true}`, "testKey", "true"},
		{`{"testKey":false}`, "testKey", "false"},
		{`{"test Key":"blue string"}`, "test Key", "blue string"},
		{`{"testKey":"blue string"}`, "testKeyNotFound", ""},
		{`{}`, "testKeyNotFound", ""},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, err := NewJson(testCase.configString)
		assert.Nil(t, err)

		// act
		result := jsonConfig.GetString(testCase.key)

		// assert
		assert.Equal(t, testCase.expected, result)
	}
}

func TestJson_GetError(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     string
		returnError  bool
	}{
		{`{"testKey":true}`, "testKey", "true", false},
		{`{"testKey":false}`, "testKey", "false", false},
		{`{"test Key":"blue string"}`, "test Key", "blue string", false},
		{`{"testKey":"blue string"}`, "testKeyNotFound", "", true},
		{`{"testKey":"blue string"}`, "testkey", "", true},
		{`{}`, "testKeyNotFound", "", true},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result, err := jsonConfig.GetStringError(testCase.key)

		// assert
		assert.Equal(t, testCase.expected, result)
		assert.Equal(t, err != nil, testCase.returnError)
	}
}

func TestJson_GetStringFatal(t *testing.T) {
	testCases := []struct {
		configString   string
		key            string
		expectedResult string
		fatalMessage   string
	}{
		{`{"testKey":true}`, "testKey", "true", ""},
		{`{"testKey":false}`, "testKey", "false", ""},
		{`{"test Key":"blue string"}`, "test Key", "blue string", ""},
		{`{"testKey":"blue string"}`, "testKeyNotFound", "<nil>", "Cannot not find configuration for testKeyNotFound"},
		{`{}`, "testKeyNotFound", "<nil>", "Cannot not find configuration for testKeyNotFound"},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, err := NewJson(testCase.configString)
		assert.Nil(t, err)

		var fatalMessage string
		jsonConfig.fatal = func(i interface{}) {
			fatalMessage = fmt.Sprint(i)
		}

		// act
		result := jsonConfig.GetStringFatal(testCase.key)

		// assert
		assert.Equal(t, testCase.expectedResult, result)
		assert.Equal(t, testCase.fatalMessage, fatalMessage)
	}
}

func TestJson_GetInt(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     int
		error        error
	}{
		{`{"testKey":true}`, "testKey", 0, &strconv.NumError{}},
		{`{"testKey":false}`, "testKey", 0, &strconv.NumError{}},
		{`{"test Key":"23"}`, "test Key", 23, nil},
		{`{"test Key":23}`, "test Key", 23, nil},
		{`{"test Key":"alpha beta"}`, "test Key", 0, &strconv.NumError{}},
		{`{"testKey":"32"}`, "testKeyNotFound", 0, errors.New("could not find testKeyNotFound")},
		{`{}`, "testKeyNotFound", 0, errors.New("could not find testKeyNotFound")},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result, err := jsonConfig.GetInt(testCase.key)

		// assert
		assert.Equal(t, testCase.expected, result)
		assert.IsType(t, testCase.error, err)
	}
}

func TestJson_GetIntFatal(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     int
		fatalMessage string
	}{
		{`{"testKey":true}`, "testKey", 0, "Cannot not find int configuration for testKey"},
		{`{"testKey":false}`, "testKey", 0, "Cannot not find int configuration for testKey"},
		{`{"test Key":"23"}`, "test Key", 23, ""},
		{`{"test Key":23}`, "test Key", 23, ""},
		{`{"test Key":"alpha beta"}`, "test Key", 0, "Cannot not find int configuration for test Key"},
		{`{"testKey":"32"}`, "testKeyNotFound", 0, "Cannot not find int configuration for testKeyNotFound"},
		{`{}`, "testKeyNotFound", 0, "Cannot not find int configuration for testKeyNotFound"},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		var fatalMessage string
		jsonConfig.fatal = func(i interface{}) {
			fatalMessage = fmt.Sprint(i)
		}

		// act
		result := jsonConfig.GetIntFatal(testCase.key)

		// assert
		assert.Equal(t, testCase.expected, result)
		assert.IsType(t, testCase.fatalMessage, fatalMessage)
	}
}

func TestJson_GetFloat(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     float64
		error        error
	}{
		{`{"testKey":true}`, "testKey", 0, &strconv.NumError{}},
		{`{"testKey":false}`, "testKey", 0, &strconv.NumError{}},
		{`{"test Key":"23.43"}`, "test Key", 23.43, nil},
		{`{"test Key":23.23}`, "test Key", 23.23, nil},
		{`{"test Key":32.3E-2}`, "test Key", 32.3e-2, nil},
		{`{"test Key":"alpha beta"}`, "test Key", 0, &strconv.NumError{}},
		{`{"testKey":"32.3E-2"}`, "testKeyNotFound", 0, errors.New("")},
		{`{}`, "testKeyNotFound", 0, errors.New("")},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result, err := jsonConfig.GetFloat(testCase.key)

		// assert
		assert.Equal(t, testCase.expected, result)
		assert.IsType(t, testCase.error, err)
	}
}

func TestJson_GetKeys(t *testing.T) {
	testCases := []struct {
		configString string
		expected     []string
	}{
		{`{"testKey":true}`, []string{"testKey"}},
		{`{"testKey1":true, "test key2":true}`, []string{"testKey1", "test key2"}},
		{`{}`, []string{}},
		{`{"testKey1":true, "test key2":true, "testKey3":{"testkey4":4, "testkey5":5}}`, []string{"testKey1", "test key2", "testKey3"}},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result := jsonConfig.GetKeys()

		// assert
		sort.Strings(testCase.expected)
		sort.Strings(result)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestJson_GetConfig(t *testing.T) {
	testCases := []struct {
		configString   string
		key            string
		expectedConfig string
	}{
		{`{"testKey":true}`, "testKey", ""},
		{`{"testKey1":true, "test key2":true}`, "testKey1", ""},
		{`{}`, "testKey", ""},
		{`{"testKey1":true, "test key2":true, "testKey3":{"testkey4":4, "testkey5":5}}`, "testKey3", `{"testkey4":4, "testkey5":5}`},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result := jsonConfig.GetConfig(testCase.key)

		// assert
		var expected Configuration
		if len(testCase.expectedConfig) > 0 {
			expected, _ = NewJson(testCase.expectedConfig)
		}
		if result == nil {
			assert.Equal(t, expected, result)
		} else {
			expectedJson, ok := expected.(*Json)
			assert.True(t, ok)
			resultJson, ok := result.(*Json)
			assert.True(t, ok)
			assert.Equal(t, expectedJson.params, resultJson.params)
		}
	}
}

func TestJson_GetStruct(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     interface{}
		error        error
	}{
		{`{"testKey":true}`, "testKey", true, nil},
		{`{"testKey":true}`, "testKey1", nil, errors.New("")},
		{`{"testKey1":true, "test key2":true, "testKey3":{"testkey4":4.5, "testkey5":5.5}}`, "testKey3", map[string]interface{}{"testkey4": 4.5, "testkey5": 5.5}, nil},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		var result interface{}
		err := jsonConfig.GetStruct(testCase.key, &result)

		// assert
		assert.IsType(t, testCase.error, err)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestJson_GetStringArray(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		expected     []string
	}{
		{`{"testKey":true}`, "testKey", nil},
		{`{"testKey":[true, false]}`, "testKey", []string{"true", "false"}},
		{`{"testKey":["hi there", "bye there"]}`, "testKey", []string{"hi there", "bye there"}},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result := jsonConfig.GetStringArray(testCase.key)

		// assert
		assert.Equal(t, testCase.expected, result)
	}
}

func TestJson_GetBool(t *testing.T) {
	testCases := []struct {
		configString string
		key          string
		defaultVal   bool
		expected     bool
	}{
		{`{"testKey":true}`, "testKey", false, true},
		{`{"testKey":false}`, "testKey", true, false},
		{`{}`, "testKey", false, false},
		{`{}`, "testKey", true, true},
		{`{"testKey":"T"}`, "testKey", false, true},
		{`{"testKey":"F"}`, "testKey", true, false},
		{`{"testKey":"True"}`, "testKey", false, true},
		{`{"testKey":"False"}`, "testKey", true, false},
		{`{"testKey":"TRUE"}`, "testKey", false, true},
		{`{"testKey":"FALSE"}`, "testKey", true, false},
		{`{"testKey":"1"}`, "testKey", false, true},
		{`{"testKey":"0"}`, "testKey", true, false},
		{`{"testKey":1}`, "testKey", false, true},
		{`{"testKey":0}`, "testKey", true, false},
	}

	// With no environmental variables
	for _, testCase := range testCases {
		// arrange
		jsonConfig, jsonErr := NewJson(testCase.configString)
		assert.Nil(t, jsonErr)

		// act
		result := jsonConfig.GetBool(testCase.key, testCase.defaultVal)

		// assert
		assert.Equal(t, testCase.expected, result, testCase.configString)
	}
}
