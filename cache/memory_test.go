// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reaction-eng/restlib/mocks"

	"github.com/golang/mock/gomock"
)

func TestNewMemory(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMemoryCache := mocks.NewMockRawMemoryCache(mockCtrl)

	// act
	memory := NewMemory(mockMemoryCache)

	// assert
	assert.Equal(t, mockMemoryCache, memory.rawMemoryCache)
}

func TestMemory_Set(t *testing.T) {
	testCases := []struct {
		key  string
		item interface{}
	}{
		{"test 1", make([]int, 3)},
		{"test 1", "blue"},
		{"test 1", 23},
		{"test 1", &struct{}{}},
	}

	mockCtrl := gomock.NewController(t)
	for _, testCase := range testCases {
		// arrange
		mockMemoryCache := mocks.NewMockRawMemoryCache(mockCtrl)
		mockMemoryCache.EXPECT().SetDefault(testCase.key, testCase.item).Times(1)

		memory := NewMemory(mockMemoryCache)

		// act
		memory.Set(testCase.key, testCase.item)
	}

	// assert
	mockCtrl.Finish()
}

func TestMemory_SetString(t *testing.T) {
	testCases := []struct {
		key  string
		item string
	}{
		{"test 1", "blue"},
		{"test 1", "blue green"},
	}

	mockCtrl := gomock.NewController(t)
	for _, testCase := range testCases {
		// arrange
		mockMemoryCache := mocks.NewMockRawMemoryCache(mockCtrl)
		mockMemoryCache.EXPECT().SetDefault(testCase.key, testCase.item).Times(1)

		memory := NewMemory(mockMemoryCache)

		// act
		memory.SetString(testCase.key, testCase.item)
	}

	// assert
	mockCtrl.Finish()
}

func TestMemory_Get(t *testing.T) {
	testCases := []struct {
		key      string
		get      interface{}
		found    bool
		expected *struct {
			Value1 int
			Value2 string
			Value3 interface{}
		}
	}{
		{"test 1", struct {
			Value1 int
			Value2 string
			Value3 interface{}
		}{Value1: 21, Value2: "test", Value3: "test"},
			true, &struct {
				Value1 int
				Value2 string
				Value3 interface{}
			}{Value1: 21, Value2: "test", Value3: "test"}},
		{"test 1", "blue", false, nil},
	}

	mockCtrl := gomock.NewController(t)
	for _, testCase := range testCases {
		// arrange
		mockMemoryCache := mocks.NewMockRawMemoryCache(mockCtrl)
		mockMemoryCache.EXPECT().Get(testCase.key).Return(testCase.get, testCase.found).Times(1)

		memory := NewMemory(mockMemoryCache)

		// act
		var result struct {
			Value1 int
			Value2 string
			Value3 interface{}
		}

		found := memory.Get(testCase.key, &result)

		assert.Equal(t, testCase.found, found)
		if found {
			assert.Equal(t, testCase.expected, &result)
		}
	}

	// assert
	mockCtrl.Finish()
}

func TestMemory_GetString(t *testing.T) {
	testCases := []struct {
		key      string
		get      string
		found    bool
		expected string
	}{
		{"testKey", "", true, ""},
		{"testKey", "123", true, "123"},
		{"testKey", "123", false, ""},
	}

	mockCtrl := gomock.NewController(t)
	for _, testCase := range testCases {
		// arrange
		mockMemoryCache := mocks.NewMockRawMemoryCache(mockCtrl)
		mockMemoryCache.EXPECT().Get(testCase.key).Return(testCase.get, testCase.found).Times(1)

		memory := NewMemory(mockMemoryCache)

		// act
		result, found := memory.GetString(testCase.key)

		assert.Equal(t, testCase.found, found)
		assert.Equal(t, testCase.expected, result)
	}

	// assert
	mockCtrl.Finish()
}
