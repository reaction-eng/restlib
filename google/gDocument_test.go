// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package google

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGDocument_GetId(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	// act
	testId := gDocument.GetId()

	//assert
	assert.Equal(t, "12345abc", testId)
}

func TestGDocument_GetName(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	// act
	testName := gDocument.GetName()

	//assert
	assert.Equal(t, "testName", testName)
}

func TestGDocument_GetDate(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
	}

	// act
	testDate := gDocument.GetDate()

	//assert
	assert.Equal(t, &referenceTime, testDate)
}

func TestGDocument_GetType(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type: "testFileType",
	}

	// act
	testFileType := gDocument.GetType()

	//assert
	assert.Equal(t, "testFileType", testFileType)
}

func TestGDocument_GetPreview(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Preview: "testFilePreview",
	}

	// act
	testPreview := gDocument.GetPreview()

	//assert
	assert.Equal(t, "testFilePreview", testPreview)
}

func TestGDocument_GetThumbnailUrl(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		ThumbnailUrl: "thumbnailUrl",
	}

	// act
	testThumbnailUrl := gDocument.GetThumbnailUrl()

	//assert
	assert.Equal(t, "thumbnailUrl", testThumbnailUrl)
}

func TestGDocument_GetParentId(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
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
	testParentId := gDocument.GetParentId()

	//assert
	assert.Equal(t, "parentId123", testParentId)
}

func TestGDocument_MarshalJSON(t *testing.T) {
	// arrange
	referenceTime := time.Unix(1574622126, 0)

	gDocument := gDocument{
		gFile: gFile{
			Id:          "12345abc",
			Name:        "testName",
			HideListing: false,
			Date:        &referenceTime,
		},
		Type:         "testFileType",
		Preview:      "testFilePreview",
		ThumbnailUrl: "utl//url",
		ParentId:     "IHaveNoFather",
	}

	// act
	jsonGDocument, _ := json.Marshal(gDocument)

	//assert
	assert.Equal(t, `{"Id":"12345abc","name":"testName","hideListing":false,"date":"2019-11-24T12:02:06-07:00","type":"testFileType","preview":"testFilePreview","thumbnail":"utl//url","parentid":"IHaveNoFather","InternalItemType":"gDocument"}`, string(jsonGDocument))
}
