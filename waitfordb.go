package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
	"time"
)

var db *sql.DB

type DBConfig struct {
	host     string
	port     int
	user     string
	password string
	name     string
}

func (dbconf *DBConfig) ParseCommandLine() {
	paramHostPtr := flag.String("host", "localhost", "Host name of the server with a MS SQLServer")
	paramPortPtr := flag.Int("port", 1433, "MSSQLServer is listing on this port")
	paramUserPtr := flag.String("user", "", "Database user name")
	paramPasswordPtr := flag.String("password", "", "Passwort of database user")
	paramNamePtr := flag.String("name", "", "Database name")

	flag.Parse()

	dbconf.host = *paramHostPtr
	dbconf.port = *paramPortPtr
	dbconf.user = *paramUserPtr
	dbconf.password = *paramPasswordPtr
	dbconf.name = *paramNamePtr

	if *paramHostPtr == "" {
		fmt.Fprintln(os.Stderr, "Parameter 'host' is empty.")
		flag.CommandLine.Usage()
		os.Exit(101)
	}
	if *paramPortPtr == 0 {
		fmt.Fprintln(os.Stderr, "Parameter 'port' is empty.")
		flag.CommandLine.Usage()
		os.Exit(102)
	}
	if *paramUserPtr == "" {
		fmt.Fprintln(os.Stderr, "Parameter 'user' is empty.")
		flag.CommandLine.Usage()
		os.Exit(103)
	}
	if *paramPasswordPtr == "" {
		dbconf.password = os.Getenv("DB_USER_PASSWORD")
		if dbconf.password == "" {
			fmt.Fprintln(os.Stderr, "Parameter user 'password' is empty. You can set also the environment 'DB_USER_PASSWORD'")
			flag.CommandLine.Usage()
			os.Exit(104)
		}
	}
	if *paramNamePtr == "" {
		fmt.Fprintln(os.Stderr, "Parameter database 'name' is empty.")
		flag.CommandLine.Usage()
		os.Exit(105)
	}
}

func main() {

	dbconf := &DBConfig{}
	dbconf.ParseCommandLine()

	// Build connection string
	connString := fmt.Sprintf("host=%s;user id=%s;password=%s;port=%d;database=%s;",
		dbconf.host, dbconf.user, dbconf.password, dbconf.port, dbconf.name)

	var tries = 0
	var triesConf = 10

	for tries < triesConf {
		var err error

		// Create connection pool
		db, err = sql.Open("sqlserver", connString)
		if err != nil {
			log.Fatalf("Error creating connection pool: ", err.Error())
		}

		ctx := context.Background()
		err = db.PingContext(ctx)

		if err == nil {
			fmt.Printf("Connected!\n")
			count, err := GetTablesCount()
			if err == nil {
				DBFound(count)
			}
		}

		switch {
		case tries == 0:
			log.Printf("Wait for Database '%s'.", dbconf.host)
		case tries < triesConf:
			log.Printf("Wait %d seconds for Database '%s'. Will give up in %d seconds.", tries*20, dbconf.name, (triesConf-tries)*20)
		}

		tries += 1
		time.Sleep(20 * time.Second)
	}
	log.Fatalf("No database '%s' found in %d seconds. Please check your configuration [host: %s, database: %s, port: %d, user: %s]",
		dbconf.name, triesConf*20, dbconf.host, dbconf.name, dbconf.port, dbconf.user)
	os.Exit(10)
}

func DBFound(tableCount int) {
	fmt.Printf("Read %d row(s) successfully.\n", tableCount)
	if tableCount == 0 {
		fmt.Printf("Database is empty! Initialization is necessary.")
		os.Exit(1)
	}
	os.Exit(0)
}

func GetTablesCount() (int, error) {
	ctx := context.Background()

	// Check if database is alive.
	err := db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := fmt.Sprintf("SELECT Distinct TABLE_NAME FROM INFORMATION_SCHEMA.TABLES;")

	// Execute query
	rows, err := db.QueryContext(ctx, tsql)
	if err != nil {
		return -1, err
	}

	defer rows.Close()

	var count = 0

	// Iterate through the result set.
	for rows.Next() {
		var tablename string

		// Get values from row.
		err := rows.Scan(&tablename)
		if err != nil {
			return -1, err
		}
		count++
	}

	return count, nil
}
