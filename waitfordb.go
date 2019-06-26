/*
 * Copyright (c) 2019.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"math/rand"
	"os"
	"time"
)

var db *sql.DB

type Config struct {
	host     string
	port     int
	user     string
	password string
	name     string
	lockfile string

	timeout    int
	timeperiod int
}

func (dbconf *Config) ParseCommandLine() {
	paramHostPtr := flag.String("host", "localhost", "Host name of the server with a MS SQLServer")
	paramPortPtr := flag.Int("port", 1433, "MSSQLServer is listing on this port")
	paramUserPtr := flag.String("user", "", "Database user name")
	paramPasswordPtr := flag.String("password", "", "Passwort of database user")
	paramNamePtr := flag.String("name", "", "Database name")
	paramLockFilePtr := flag.String("lockfile", "", "Create a lock file")

	paramTimeout := flag.Int("timeout", 200, "Timeout for waiting in seconds")
	paramTimeperiod := flag.Int("timeperiod", 20, "Time between checks in seconds")

	flag.Parse()

	dbconf.host = *paramHostPtr
	dbconf.port = *paramPortPtr
	dbconf.user = *paramUserPtr
	dbconf.password = *paramPasswordPtr
	dbconf.name = *paramNamePtr
	//locking table
	dbconf.lockfile = *paramLockFilePtr
	//time configuration
	dbconf.timeout = *paramTimeout
	dbconf.timeperiod = *paramTimeperiod

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
	if *paramTimeout <= *paramTimeperiod {
		fmt.Fprintln(os.Stderr, "Parameter timeperiod is bigger or equal to timeout. The timeout configuration must be bigger than the timeperiod.")
		flag.CommandLine.Usage()
		os.Exit(106)
	}
	if *paramTimeout <= 0 {
		fmt.Fprintln(os.Stderr, "Parameter timeout must be bigger than 0.")
		flag.CommandLine.Usage()
		os.Exit(107)
	}
	if *paramTimeperiod <= 0 {
		fmt.Fprintln(os.Stderr, "Parameter timeperiod must be bigger than 0.")
		flag.CommandLine.Usage()
		os.Exit(108)
	}
}

func main() {

	dbconf := &Config{}
	dbconf.ParseCommandLine()

	// Build connection string
	connString := fmt.Sprintf("host=%s;user id=%s;password=%s;port=%d;database=%s;",
		dbconf.host, dbconf.user, dbconf.password, dbconf.port, dbconf.name)

	runTime := 0
	available := false

	for runTime < dbconf.timeout {
		var err error

		// check of lock file
		if dbconf.lockfile != "" {
			available = !dbconf.LockFileExists()
		}

		// check database
		if (dbconf.lockfile != "" && available) || dbconf.lockfile == "" {
			// Create connection pool
			db, err = sql.Open("sqlserver", connString)
			if err != nil {
				log.Fatalf("Error creating connection pool: %s", err.Error())
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
		}

		log.Printf("Wait %d seconds for Database '%s'. Will give up in %d seconds.", runTime, dbconf.name, (dbconf.timeout - runTime))

		runTime += dbconf.timeperiod
		time.Sleep(time.Duration(rand.Int31n(int32(dbconf.timeperiod))) * time.Second)
	}

	log.Fatalf("No database '%s' found in %d seconds. Please check your configuration [host: %s, database: %s, port: %d, user: %s]",
		dbconf.name, runTime, dbconf.host, dbconf.name, dbconf.port, dbconf.user)
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

func (dbconf *Config) LockFileExists() bool {
	if _, err := os.Stat(dbconf.lockfile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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
