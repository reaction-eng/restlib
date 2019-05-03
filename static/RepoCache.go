package static

import (
	"bitbucket.org/reidev/restlib/cache"
	"bitbucket.org/reidev/restlib/configuration"
	"bitbucket.org/reidev/restlib/google"
)

/**
Define a struct for RepoMem for news
*/
type RepoGoogleCache struct {
	//Store the cache
	cas cache.ObjectCache

	//We also need googl
	drive *google.Drive

	//Store the public and private
	privateConfig *configuration.Configuration
	publicConfig  *configuration.Configuration
}

//Provide a method to make a new AnimalRepoSql
func NewRepoCache(drive *google.Drive, cas cache.ObjectCache, privateConfigFile string, publicConfigFile string) *RepoGoogleCache {

	//Create a new config
	privateConfig, _ := configuration.NewConfiguration(privateConfigFile)
	publicConfig, _ := configuration.NewConfiguration(publicConfigFile)

	//Define a new repo
	newRepo := RepoGoogleCache{
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
func (repo *RepoGoogleCache) GetStaticPublicDocument(path string) (string, error) {

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
func (repo *RepoGoogleCache) GetStaticPrivateDocument(path string) (string, error) {

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
func (repo *RepoGoogleCache) CleanUp() {

}
