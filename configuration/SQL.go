package configuration

import "fmt"

//Build the dbString //username:password@protocol(address)/dbname
func (config *Configuration) GetMySqlDataBaseSourceName() string {
	dbString := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true",
		config.GetString("db_username"),
		config.GetString("db_password"),
		config.GetString("db_protocol"),
		config.GetString("db_address"),
		config.GetString("db_name"),
	)

	return dbString ////"root:P1p3sh0p@tcp(:3306)/localDB?parseTime=true"
}

//Build the dbString //username:password@protocol(address)/dbname
func (config *Configuration) GetPostgresDataBaseSourceName() string {
	//dbString :=   "postgres://postgres:kOVGMnoS3iIk@localhost/postgres?sslmode=disable"
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		config.GetString("db_username"),
		config.GetString("db_password"),
		config.GetString("db_address"),
		config.GetString("db_name"),
	)

	return dbString
}
