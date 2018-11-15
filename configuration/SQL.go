package configuration

import "fmt"

//Build the dbString //username:password@protocol(address)/dbname
func (config *Configuration) GetDataBaseSourceName() string {
	dbString := fmt.Sprintf("%s:%s@%s(%s)/%s",
		config.GetString("db_username"),
		config.GetString("db_password"),
		config.GetString("db_protocol"),
		config.GetString("db_address"),
		config.GetString("db_name"),
	)

	return dbString
}
