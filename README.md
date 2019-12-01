# REI RESTLib #

## Introduction ##

REI RESTLib is a Go library providing HTTP REST APIs for various features useful for building the backend for a single page web application (SPA). These components include users, roles, several authentication methods, user preferences, and more.

The code is still under development and currently makes no guarantees of API stability.

## Usage ##

TODO

## License ##

This library was originally developed by Reaction Engineering International for several internal projects. It is now being made available under the MIT license (see LICENSE.txt).

## Mocking ##
You can update or generate all mocks for testing with:
```
go generate ./...
```


## Testing ##
The Framework uses the builtin testing capabilities in go

## DB Migrations ##
All db migrations are handled using sql-migrate found at https://github.com/rubenv/sql-migrate. When using this as a library, you can act upon the migrations
```go
	db, err := sql.Open("postgres", dbString) //"root:P1p3sh0p@tcp(:3306)/localDB?parseTime=true"
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	restLibMigrations := migrations.Postgres()

	migrate.SetTable("restlib_migrations")
	n, err := migrate.Exec(db, "postgres", restLibMigrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)

```