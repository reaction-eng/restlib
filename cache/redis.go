// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
)

type Redis struct {
	codec *cache.Codec
}

func NewObjectRedisCache(redis *redis.Ring) *Redis {
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
	infoRepo := Redis{
		codec: codec,
	}

	//Return a point
	return &infoRepo
}

func (repo *Redis) Get(key string, item interface{}) bool {
	//Get the summary
	err := repo.codec.Get(key, &item)
	if err != nil {
		item = nil
		return false
	}
	return true
}

func (repo *Redis) Set(key string, item interface{}) error {
	return repo.codec.Set(&cache.Item{
		Key:        key,
		Object:     item,
		Expiration: time.Hour,
	})
}

func (repo *Redis) GetString(key string) (string, bool) {
	var item string

	err := repo.codec.Get(key, &item)
	if err != nil {
		return "", false
	}

	return item, true
}

func (repo *Redis) SetString(key string, value string) {
	repo.codec.Set(&cache.Item{
		Key:        key,
		Object:     value,
		Expiration: time.Hour,
	})
}
