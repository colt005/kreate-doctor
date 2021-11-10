package dbcontroller

import (
	"database/sql"
	"fmt"
	"os"
)

var DB *sql.DB

func PreparePSQL() (db *sql.DB, err error) {

	var host string = os.Getenv("PGHOST")
	var port string = os.Getenv("PGPORT")
	var user string = os.Getenv("PGUSER")
	var password string = os.Getenv("PGPASSWORD")
	var dbname string = os.Getenv("PGDATABASE")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	errp := database.Ping()

	if errp != nil {
		fmt.Println(errp.Error())
	} else {
		fmt.Println("Successfully pinged!")
	}
	fmt.Println("Successfully connected!")

	DB = database
	return database, nil
}
