package timescale

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

type TimeScale struct {
	connString string
	connection *sql.DB
}

var db TimeScale

func SetupDb() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db = TimeScale{connString: os.Getenv("DATABASE_OPTS")}
	db.connection, err = sql.Open("postgres", db.connString)

	if err != nil {
		log.Fatal(err)
	}
}

func InsertTradeData(price float32, amount float32, pair string, tableName string) {

	query := []string{"INSERT INTO ", tableName, "(time, price, amount, pair) VALUES (NOW(), $1, $2, $3);"}
	_, err := db.connection.Exec(strings.Join(query, ""), price, amount, pair)

	if err != nil {
		log.Fatal(err)
	}
}
