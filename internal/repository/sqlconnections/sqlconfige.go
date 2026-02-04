package sqlconnections

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDb() (*sql.DB, error) {

	err := godotenv.Load("cmd/api/.env")
	if err != nil {
		log.Println("Error in loading .env content")
		return nil, err
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWOrd")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("HOST")
	port := os.Getenv("DB_PORT")

	connection := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", user, password, host, port, dbname)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		// panic(err)
		return nil, err
	}
	
	return db, nil
}
