// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package cache

import (
	"encoding/json"
	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"time"
)

/**
Define a struct for RepoMem for news
*/
type ObjectRedisCache struct {
	//Store the cache
	codec *cache.Codec
}

//Provide a method to make a new AnimalRepoSql
func NewObjectRedisCache(redis *redis.Ring) *ObjectRedisCache {

	codec := &cache.Codec{
		Redis: redis,

		Marshal: func(v interface{}) ([]byte, error) {
			return json.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return json.Unmarshal(b, v)
		},
	}

	//Define a new repo
	infoRepo := ObjectRedisCache{
		codec: codec,
	}

	//Return a point
	return &infoRepo

}

/**
Get all of the news
*/
func (repo *ObjectRedisCache) Get(key string, item interface{}) {

	//Get the summary
	err := repo.codec.Get(key, &item)
	if err != nil {
		item = nil
	}

}

/**
Get all of the news
*/
func (repo *ObjectRedisCache) Set(key string, item interface{}) error {

	//Now save it
	return repo.codec.Set(&cache.Item{
		Key:        key,
		Object:     item,
		Expiration: time.Hour,
	})

}

func (repo *ObjectRedisCache) GetString(key string) (string, bool) {

	//Get the google item
	var item string

	//Get the summary
	err := repo.codec.Get(key, &item)
	if err != nil {
		return "", false
	}
	//
	////Now return the item
	return item, true

}

/**
Get all of the news
*/
func (repo *ObjectRedisCache) SetString(key string, value string) {

	//Now save it
	repo.codec.Set(&cache.Item{
		Key:        key,
		Object:     value,
		Expiration: time.Hour,
	})

}
