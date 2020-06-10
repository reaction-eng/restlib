// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package static

import (
	"github.com/reaction-eng/restlib/cache"
	"github.com/reaction-eng/restlib/configuration"
	"github.com/reaction-eng/restlib/file"
)

type RepoCache struct {
	//Store the cache
	cas cache.Cache

	drive file.Storage

	//Store the public and private
	privateConfig configuration.Configuration
	publicConfig  configuration.Configuration
}

func NewRepoCache(drive file.Storage, cas cache.Cache, privateConfig configuration.Configuration, publicConfig configuration.Configuration) *RepoCache {

	newRepo := RepoCache{
		cas:           cas,
		drive:         drive,
		privateConfig: privateConfig,
		publicConfig:  publicConfig,
	}

	//Return a point
	return &newRepo

}

/**
Get the public static
*/
func (repo *RepoCache) GetStaticPublicDocument(path string) (string, error) {

	//Look up the document id from the config
	documentId, err := repo.publicConfig.GetStringError(path)

	if err != nil {
		return "", err
	}

	//see if there is a cache, if there is no cache just return it
	if repo.cas == nil {
		return repo.drive.GetFileHtml(documentId), nil
	}

	//Get the summary
	value, found := repo.cas.GetString(documentId)
	if !found {
		//Update it
		value = repo.drive.GetFileHtml(documentId)

		//Now save it
		repo.cas.SetString(documentId, value)
	}

	//Now return the item
	return value, nil

}

/**
Get the public static
*/
func (repo *RepoCache) GetStaticPrivateDocument(path string) (string, error) {
	//Look up the document id from the config
	documentId, err := repo.privateConfig.GetStringError(path)

	if err != nil {
		return "", err
	}

	//see if there is a cache, if there is no cache just return it
	if repo.cas == nil {
		return repo.drive.GetFileHtml(documentId), nil
	}

	//Get the summary
	value, found := repo.cas.GetString(documentId)
	if !found {
		//Update it
		value = repo.drive.GetFileHtml(documentId)

		//Now save it
		repo.cas.SetString(documentId, value)
	}

	//Now return the item
	return value, nil

}

/**
Nothing much to do for the clean up
*/
func (repo *RepoCache) CleanUp() {

}
