package cache

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"

	"time"
)

/**
Define a struct for RepoMem for news
*/
type ObjectMemCache struct {
	//Store the cache
	cache *cache.Cache
}

//Provide a method to make a new AnimalRepoSql
func NewObjectMemCache() *ObjectMemCache {

	//Define a new repo
	infoRepo := ObjectMemCache{
		cache: cache.New(time.Hour/2.0, time.Hour),
	}

	//Return a point
	return &infoRepo

}

/**
Get all of the news
*/
func (repo *ObjectMemCache) Get(key string, returnItem interface{}) {

	item, found := repo.cache.Get(key)

	if found {

		//Convert to json
		jsonByte, err := json.Marshal(item)

		//If there is no error
		if err != nil {
			return
		}

		//Now restore back
		json.Unmarshal(jsonByte, &returnItem)

	} else {
		returnItem = nil
	}

}

/**
Get all of the news
*/
func (repo *ObjectMemCache) Set(key string, item interface{}) error {

	//Now save it
	repo.cache.SetDefault(key, item)
	return nil

}

func (repo *ObjectMemCache) GetString(key string) (string, bool) {

	item, found := repo.cache.Get(key)

	if found {
		return item.(string), true
	} else {
		return "", false
	}

}

/**
Get all of the news
*/
func (repo *ObjectMemCache) SetString(key string, value string) {

	//Now save it
	repo.cache.SetDefault(key, value)

}
