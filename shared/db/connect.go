package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"shared/enum"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS data (id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL, email TEXT NOT NULL, phone TEXT NOT NULL, company TEXT NOT NULL, ccn TEXT NOT NULL, designation TEXT NOT NULL)")
	if err != nil {
		log.Fatal("Cannot create table.", err)
	}
}

func DeleteTable(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS data CASCADE;")
	if err != nil {
		log.Fatal("Cannot drop table.", err)
	}
}

func ConnectDB() *sql.DB {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pghost := os.Getenv("POSTGRES_HOST")
	pgport := os.Getenv("POSTGRES_PORT")
	pguser := os.Getenv("POSTGRES_USER")
	pgpassword := os.Getenv("POSTGRES_PASSWORD")
	pgdbname := os.Getenv("POSTGRES_DBNAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		pghost, pgport, pguser, pgpassword, pgdbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	} else if err = db.Ping(); err != nil {
		panic(err)
	}

	createTable(db)

	return db

}

func Insert(db *sql.DB, data enum.Data) {
	sqlStatement := "INSERT INTO data (name, email, phone, company, ccn, designation) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(sqlStatement, data.Name, data.Email, data.Phone, data.Company, data.Ccn, data.Designation)
	if err != nil {
		log.Fatal("Error inserting values -->", err)
	}
}
