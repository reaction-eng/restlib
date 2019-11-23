// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package configuration

//go:generate mockgen -destination=../mocks/mock_configuration.go -package=mocks github.com/reaction-eng/restlib/configuration Configuration

type Configuration interface {
	Get(key string) interface{}
	GetFatal(key string) interface{}
	GetString(key string) string
	GetStringError(key string) (string, error)
	GetStringFatal(key string) string
	GetInt(key string) (int, error)
	GetIntFatal(key string) int
	GetFloat(key string) (float64, error)
	GetKeys() []string
	GetConfig(key string) Configuration
	GetStruct(key string, object interface{}) error
	GetStringArray(key string) []string
	GetBool(key string, defaultVal bool) bool
}
