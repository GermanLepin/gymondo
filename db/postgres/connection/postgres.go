package connection

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose"
	"log"
	"os"
	"time"

	_ "gymondo/db/postgres/migrations"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var driver = "pgx"

func StartDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s database=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DATABASE"),
	)

	conn := connectToDB(dsn)
	if conn == nil {
		return nil, fmt.Errorf("cannot connect to Postgres")
	}

	if err := goose.Up(conn, "/var"); err != nil {
		return nil, fmt.Errorf("cannot run the migrations, error is: %s", err)
	}

	// if smth goes wrong we always can run down Migrations goose.Down()
	// if err := goose.Down(conn, "/var"); err != nil {
	//    return nil, fmt.Errorf("cannot run the migrations, error is: %s", err)
	// }

	return conn, nil
}

func connectToDB(dsn string) *sql.DB {
	var counts int64

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres is not ready yet")
			counts++
		} else {
			log.Println("connected to Postgres")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
