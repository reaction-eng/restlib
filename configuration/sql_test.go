// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package configuration_test

import (
	"testing"

	"github.com/reaction-eng/restlib/configuration"

	"github.com/golang/mock/gomock"
	"github.com/reaction-eng/restlib/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewSql(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)

	// act
	sqlConfig := configuration.NewSql(mockConfiguration)

	//assert
	assert.Equal(t, mockConfiguration, sqlConfig.Configuration)
}

func TestSql_GetMySqlDataBaseSourceName(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetString("db_username").Times(1).Return("DBUSERNAME")
	mockConfiguration.EXPECT().GetString("db_password").Times(1).Return("DBPASSWORD")
	mockConfiguration.EXPECT().GetString("db_protocol").Times(1).Return("DBPROTOCOL")
	mockConfiguration.EXPECT().GetString("db_address").Times(1).Return("DBADDRESS")
	mockConfiguration.EXPECT().GetString("db_name").Times(1).Return("DBNAME")

	sqlConfig := configuration.NewSql(mockConfiguration)

	// act
	dbString := sqlConfig.GetMySqlDataBaseSourceName()

	//assert
	assert.Equal(t, "DBUSERNAME:DBPASSWORD@DBPROTOCOL(DBADDRESS)/DBNAME?parseTime=true", dbString)
}

func TestSql_GetPostgresDataBaseSourceName(t *testing.T) {
	// arrange
	mockCtrl := gomock.NewController(t)
	mockConfiguration := mocks.NewMockConfiguration(mockCtrl)
	mockConfiguration.EXPECT().GetString("db_username").Times(1).Return("DBUSERNAME")
	mockConfiguration.EXPECT().GetString("db_password").Times(1).Return("DBPASSWORD")
	mockConfiguration.EXPECT().GetString("db_address").Times(1).Return("DBADDRESS")
	mockConfiguration.EXPECT().GetString("db_name").Times(1).Return("DBNAME")

	sqlConfig := configuration.NewSql(mockConfiguration)

	// act
	dbString := sqlConfig.GetPostgresDataBaseSourceName()

	//assert
	assert.Equal(t, "postgres://DBUSERNAME:DBPASSWORD@DBADDRESS/DBNAME?sslmode=disable", dbString)
}
