// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package configuration

import "fmt"

type Sql struct {
	Configuration
}

func NewSql(configuration Configuration) *Sql {
	return &Sql{
		configuration,
	}
}

//Build the dbString //username:password@protocol(address)/dbname
func (sqlConfig Sql) GetMySqlDataBaseSourceName() string {
	dbString := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true",
		sqlConfig.GetString("db_username"),
		sqlConfig.GetString("db_password"),
		sqlConfig.GetString("db_protocol"),
		sqlConfig.GetString("db_address"),
		sqlConfig.GetString("db_name"),
	)

	return dbString ////"root:P1p3sh0p@tcp(:3306)/localDB?parseTime=true"
}

//Build the dbString //username:password@protocol(address)/dbname
func (sqlConfig Sql) GetPostgresDataBaseSourceName() string {
	//dbString :=   "postgres://postgres:kOVGMnoS3iIk@localhost/postgres?sslmode=disable"
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		sqlConfig.GetString("db_username"),
		sqlConfig.GetString("db_password"),
		sqlConfig.GetString("db_address"),
		sqlConfig.GetString("db_name"),
	)

	return dbString
}
