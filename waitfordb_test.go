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
	"github.com/Flaque/filet"
	_ "github.com/Flaque/filet"
	"os"
	"testing"
)

func TestParseCommandLineDatabase(t *testing.T) {
	os.Args = []string{"command", "--jdbcurl=jdbc:oracle:thin:@hostname:1521:sid", "--user=username", "--password=passwd"}

	m := &Config{}
	m.ParseCommandLine()

	if m.timeout != 200 {
		t.Errorf("Timeout is not configured with default value. It is %d.", m.timeout)
	}
}

/**
func TestParseCommandLineLockFile(t *testing.T) {
	os.Args = []string{"command", "--jdbcurl=jdbc:oracle:thin:@hostname:1521:sid", "--user=username", "--password=passwd", "--lockfile=test"}

	mc := &Config{}
	mc.ParseCommandLine()

	if mc.lockfile != "test" {
		t.Errorf("Lockfile is not configured. It is %s.", mc.lockfile)
	}
}

func TestOracleDB(t *testing.T) {

	dbConnAvailable := -1

	dbcon := &DBConnection{}
	config := &Config{}

	config.user = "intershop"
	config.password = "intershop"

	dbcon.SetDBParamsFromJDBC("jdbc:oracle:thin:@localhost:1521:XE")
	dbcon.SetConnectionString(*config)

	dbConnAvailable = CheckOracleDB(dbcon)

	if dbConnAvailable == 1 {
		t.Errorf("DBConnection is not available")
	}
}


func TestSqlServerDB(t *testing.T) {

	dbConnAvailable := -1

	dbcon := &DBConnection{}
	config := &Config{}

	config.user = "intershop"
	config.password = "intershop"

	dbcon.SetDBParamsFromJDBC("jdbc:sqlserver://localhost:1433;databaseName=icmdb")
	dbcon.SetConnectionString(*config)

	dbConnAvailable = CheckSQLServerDB(dbcon)

	if dbConnAvailable == 0 {
		t.Errorf("DBConnection is not available")
	}
}

**/

func TestFileExist(t *testing.T) {
	defer filet.CleanUp(t)

	// Creates a temporary file with string "some content"
	file := filet.TmpFile(t, "", "some content")

	md := &Config{}
	md.lockfile = file.Name()

	if !md.LockFileExists() {
		t.Errorf("Lockfile is there but not correct identified. File is %s", file.Name())
	}
}

func TestOracleParameter(t *testing.T) {
	dbcon := &DBConnection{}

	dbcon.SetDBParamsFromJDBC("jdbc:oracle:thin:@hostname:1521:sid")

	if dbcon.name != "sid" || dbcon.port != 1521 || dbcon.host != "hostname" {
		t.Errorf("DB is not correct identified!")
	}
}

func TestMSSQLParameter(t *testing.T) {
	dbcon := &DBConnection{}

	dbcon.SetDBParamsFromJDBC("jdbc:sqlserver://icm-mssql-server:1433;databaseName=icmdb")

	if dbcon.name != "icmdb" || dbcon.port != 1433 || dbcon.host != "icm-mssql-server" {
		t.Errorf("DB is not correct identified! %s, %d, %s", dbcon.name, dbcon.port, dbcon.host)
	}
}
