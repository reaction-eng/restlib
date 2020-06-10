// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package cache

//go:generate mockgen -destination=../mocks/mock_memory.go -package=mocks github.com/reaction-eng/restlib/cache RawMemoryCache

import (
	"encoding/json"
)

type Memory struct {
	rawMemoryCache RawMemoryCache
}

type RawMemoryCache interface {
	Get(k string) (interface{}, bool)
	SetDefault(k string, x interface{})
}

func NewMemory(rawMemoryCache RawMemoryCache) *Memory {
	objectMemCache := Memory{
		rawMemoryCache,
	}

	return &objectMemCache
}

func (memory *Memory) Set(key string, item interface{}) error {
	//Now save it
	memory.rawMemoryCache.SetDefault(key, item)
	return nil
}

func (memory *Memory) SetString(key string, value string) {
	//Now save it
	memory.rawMemoryCache.SetDefault(key, value)
}

func (memory *Memory) Get(key string, returnItem interface{}) bool {
	item, found := memory.rawMemoryCache.Get(key)

	if found {
		//Convert to json
		jsonByte, err := json.Marshal(item)

		if err != nil {
			return false
		}

		//Now restore back
		json.Unmarshal(jsonByte, &returnItem)
		return true
	} else {
		return false
	}
}

func (memory *Memory) GetString(key string) (string, bool) {
	item, found := memory.rawMemoryCache.Get(key)

	if found {
		return item.(string), true
	}
	return "", false
}
