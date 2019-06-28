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
	"github.com/pkg/errors"
	_ "gopkg.in/goracle.v2"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

type Config struct {
	jdbcurl  string
	user     string
	password string
	lockfile string

	timeout    int
	timeperiod int
}

type DBConnection struct {
	host string
	port int
	name string

	dbtype string
}

func (config *Config) ParseCommandLine() {
	paramJdbcUrl := flag.String("jdbcurl", "", "JDBC connection URL ")
	paramUserPtr := flag.String("user", "", "Database user name")
	paramPasswordPtr := flag.String("password", "", "Passwort of database user")
	paramLockFilePtr := flag.String("lockfile", "", "Create a lock file")

	paramTimeout := flag.Int("timeout", 200, "Timeout for waiting in seconds")
	paramTimeperiod := flag.Int("timeperiod", 20, "Time between checks in seconds")

	flag.Parse()

	config.jdbcurl = *paramJdbcUrl
	config.user = *paramUserPtr
	config.password = *paramPasswordPtr
	//locking table
	config.lockfile = *paramLockFilePtr
	//time configuration
	config.timeout = *paramTimeout
	config.timeperiod = *paramTimeperiod

	if *paramJdbcUrl == "" {
		fmt.Fprintln(os.Stderr, "Parameter 'jdbcurl' is empty.")
		flag.CommandLine.Usage()
		os.Exit(101)
	}
	if *paramUserPtr == "" {
		fmt.Fprintln(os.Stderr, "Parameter 'user' is empty.")
		flag.CommandLine.Usage()
		os.Exit(103)
	}
	if *paramPasswordPtr == "" {
		config.password = os.Getenv("DB_USER_PASSWORD")
		if config.password == "" {
			fmt.Fprintln(os.Stderr, "Parameter user 'password' is empty. You can set also the environment 'DB_USER_PASSWORD'")
			flag.CommandLine.Usage()
			os.Exit(104)
		}
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

	config := &Config{}
	config.ParseCommandLine()

	dbconfig := &DBConnection{}

	// Build connection string
	connString := fmt.Sprintf("host=%s;user id=%s;password=%s;port=%d;database=%s;",
		dbconfig.host, config.user, config.password, dbconfig.port, dbconfig.name)

	runTime := 0
	available := false

	for runTime < config.timeout {
		var err error

		// check of lock file
		if config.lockfile != "" {
			available = !config.LockFileExists()
		}

		// check database
		if (config.lockfile != "" && available) || config.lockfile == "" {
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

		log.Printf("Wait %d seconds for Database '%s'. Will give up in %d seconds.", runTime, dbconfig.name, (config.timeout - runTime))

		runTime += config.timeperiod
		time.Sleep(time.Duration(rand.Int31n(int32(config.timeperiod))) * time.Second)
	}

	log.Fatalf("No database '%s' found in %d seconds. Please check your configuration [host: %s, database: %s, port: %d, user: %s]",
		dbconfig.name, runTime, dbconfig.host, dbconfig.name, dbconfig.port, config.user)
	os.Exit(10)
}

func (dbconn *DBConnection) SetDBParamsFromJDBC(jdbcurl string) error {
	if !strings.HasPrefix(jdbcurl, "jdbc:") {
		return errors.Errorf("JDBC url parameter is not correct. It does not start with 'jdbc:' (%s)", jdbcurl)
	}
	urlParts := strings.Split(jdbcurl, ":")
	// check for oracle
	// jdbc:oracle:thin:@hostname:1521:sid
	if urlParts[1] == "oracle" {
		dbconn.dbtype = "oracle"
		if len(urlParts) > 3 {
			tempHostname := urlParts[3]

			if strings.HasPrefix(tempHostname, "@") {
				dbconn.host = tempHostname[1:]
			} else {
				return errors.Errorf("JDBC url parameter is not correct. There is no hostname. (%s)", jdbcurl)
			}
		}
		if len(urlParts) > 4 {
			dbconn.port, _ = strconv.Atoi(urlParts[4])
		}

		if len(urlParts) > 5 {
			dbconn.name = urlParts[5]
		}
	}

	// check for mssql
	// jdbc:sqlserver://icm-mssql-server:1433;databaseName=icmdb
	if urlParts[1] == "sqlserver" {
		dbconn.dbtype = "mssql"
		if len(urlParts) > 2 {
			tempHostname := urlParts[2]

			if strings.HasPrefix(tempHostname, "//") {
				dbconn.host = tempHostname[2:]
			} else {
				return errors.Errorf("JDBC url parameter is not correct. There is no hostname. (%s)", jdbcurl)
			}
		}
		if len(urlParts) > 3 {
			tempPortName := strings.Split(urlParts[3], ";")
			if len(tempPortName) == 2 {
				dbconn.port, _ = strconv.Atoi(tempPortName[0])

				tempDBName := strings.Split(tempPortName[1], "=")
				if len(tempDBName) > 1 {
					dbconn.name = tempDBName[1]
				}
			}
		}
	}

	return errors.Errorf("JDBC url parameter is not correct. This is not an oracle or ms sql url (%s)", jdbcurl)
}

func DBFound(tableCount int) {
	fmt.Printf("Read %d row(s) successfully.\n", tableCount)
	if tableCount == 0 {
		fmt.Printf("Database is empty! Initialization is necessary.")
		os.Exit(1)
	}
	os.Exit(0)
}

func (config *Config) LockFileExists() bool {
	if _, err := os.Stat(config.lockfile); err != nil {
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
