// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/reaction-eng/restlib/file"

	"github.com/stretchr/testify/assert"
)

func TestGDirectory_GetId(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	// act
	testId := gDirectory.GetId()

	//assert
	assert.Equal(t, "12345abc", testId)
}

func TestGDirectory_GetName(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	// act
	testName := gDirectory.GetName()

	//assert
	assert.Equal(t, "testName", testName)
}

func TestGDirectory_GetDate(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	// act
	testDate := gDirectory.GetDate()

	//assert
	assert.Equal(t, &referenceTime, testDate)
}

func TestGDirectory_GetType(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type: "testFileType",
	}

	// act
	testFileType := gDirectory.GetType()

	//assert
	assert.Equal(t, "testFileType", testFileType)
}

func TestGDirectory_GetItems(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	expectedItems := []file.Item{
		gFile{
			Id:          "subItemId1",
			Name:        "subItemName1",
			HideListing: false,
			Date:        &referenceTime,
		},
		gFile{
			Id:          "subItemName2",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type:  "testFileType",
		Items: expectedItems,
	}

	// act
	testItems := gDirectory.GetItems()

	//assert
	assert.Equal(t, expectedItems, testItems)
}

func TestGDirectory_GetParentId(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type:     "testFileType",
		ParentId: "parentId123",
	}

	// act
	testParentId := gDirectory.GetParentId()

	//assert
	assert.Equal(t, "parentId123", testParentId)
}

func TestGDirectory_ForEach(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	items := []file.Item{
		&gFile{
			Id:          "subItemId1",
			Name:        "subItemName1",
			HideListing: false,
			Date:        &referenceTime,
		},
		&gFile{
			Id:          "subItemName2",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		&gDirectory{
			gFile: gFile{
				Id: "subDirId1",
			},
			Type:     "",
			ParentId: "",
			Items: []file.Item{
				&gFile{
					Id:          "subSubItemId1",
					Name:        "subSubItemName1",
					HideListing: false,
					Date:        &referenceTime,
				},
			},
		},
	}

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type:  "testFileType",
		Items: items,
	}

	testCases := []struct {
		includeDir bool
		expected   []string
	}{
		{
			includeDir: false,
			expected:   []string{"subItemId1", "subItemName2", "subSubItemId1"},
		},
		{
			includeDir: true,
			expected:   []string{"12345abc", "subItemId1", "subItemName2", "subSubItemId1", "subDirId1"},
		},
	}
	for _, testCase := range testCases {
		// act
		calledCommands := make([]string, 0)
		gDirectory.ForEach(testCase.includeDir, func(item file.Item) {
			calledCommands = append(calledCommands, item.GetId())
		})

		//assert
		sort.Strings(calledCommands)
		sort.Strings(testCase.expected)
		assert.Equal(t, testCase.expected, calledCommands)
	}
}

func TestGDirectory_MarshalJSON(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	items := []file.Item{
		&gDocument{
			gFile: gFile{
				Id:          "subItemId1",
				Name:        "subItemName1",
				HideListing: false,
				Date:        &referenceTime,
			},
			Preview: "Preview123",
		},
		&gFile{
			Id:          "subItemName2",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		&gDirectory{
			gFile: gFile{
				Id: "subDirId1",
			},
			Type:     "",
			ParentId: "",
			Items: []file.Item{
				&gFile{
					Id:          "subSubItemId1",
					Name:        "subSubItemName1",
					HideListing: false,
					Date:        &referenceTime,
				},
			},
		},
	}

	gDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type:  "testFileType",
		Items: items,
	}

	// act
	jsonGDirectory, _ := json.Marshal(gDirectory)

	//assert
	assert.Equal(t, `{"Id":"12345abc","name":"testName","hideListing":false,"date":"2019-11-24T12:02:06-07:00","type":"testFileType","Items":[{"Id":"subItemId1","name":"subItemName1","hideListing":false,"date":"2019-11-24T12:02:06-07:00","type":"","preview":"Preview123","thumbnail":"","parentid":"","InternalItemType":"gDocument"},{"Id":"subItemName2","name":"testName","hideListing":false,"date":"2019-11-24T12:02:06-07:00"},{"Id":"subDirId1","name":"","hideListing":false,"date":null,"type":"","Items":[{"Id":"subSubItemId1","name":"subSubItemName1","hideListing":false,"date":"2019-11-24T12:02:06-07:00"}],"parentid":"","InternalItemType":"gDirectory"}],"parentid":"","InternalItemType":"gDirectory"}`, string(jsonGDirectory))
}

func TestGDirectory_UnmarshalJSON(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	expectedItems := []file.Item{
		&gDocument{
			gFile: gFile{
				Id:          "subItemId1",
				Name:        "subItemName1",
				HideListing: false,
				Date:        &referenceTime,
			},
			Preview: "Preview123",
		},
		&gFile{
			Id:          "subItemName2",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		&gDirectory{
			gFile: gFile{
				Id: "subDirId1",
			},
			Type:     "",
			ParentId: "",
			Items: []file.Item{
				&gFile{
					Id:          "subSubItemId1",
					Name:        "subSubItemName1",
					HideListing: false,
					Date:        &referenceTime,
				},
			},
		},
	}

	expectedGDirectory := gDirectory{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type:  "testFileType",
		Items: expectedItems,
	}

	jsonString := `{"Id":"12345abc","name":"testName","hideListing":false,"date":"2019-11-24T12:02:06-07:00","type":"testFileType","Items":[{"Id":"subItemId1","name":"subItemName1","hideListing":false,"date":"2019-11-24T12:02:06-07:00","type":"","preview":"Preview123","thumbnail":"","parentid":"","InternalItemType":"gDocument"},{"Id":"subItemName2","name":"testName","hideListing":false,"date":"2019-11-24T12:02:06-07:00"},{"Id":"subDirId1","name":"","hideListing":false,"date":null,"type":"","Items":[{"Id":"subSubItemId1","name":"subSubItemName1","hideListing":false,"date":"2019-11-24T12:02:06-07:00"}],"parentid":"","InternalItemType":"gDirectory"}],"parentid":"","InternalItemType":"gDirectory"}`

	// act
	var testGDirectory gDirectory
	err := json.Unmarshal([]byte(jsonString), &testGDirectory)

	//assert
	assert.Nil(t, err)
	assert.EqualValues(t, expectedGDirectory, testGDirectory)
}
