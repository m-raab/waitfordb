Wait For DB (MS SQL and Oracle)
=========================
This tool checks the configured database connection. If the database is not available the command exits with a special exit code.

Prerequisites
-------------------------
And Oracle instant client must be installed:
See https://www.oracle.com/database/technologies/instant-client.html

Download
- Linux: instantclient-basiclite-linux.x64-12.2.0.1.0.zip
- MacOS: instantclient-basiclite-macos.x64-12.2.0.1.0.zip

Usage
-------------------------

It is possible to configure the following parameters:

```
Usage of waitfordb:
  -jdbcurl string
        JDBC url for Oracle or SQLServer
        eg. jdbc:oracle:thin:@localhost:1521:XE or 
        or jdbc:sqlserver://localhost:1433;databaseName=icmdb
  -password string
    	Passwort of database user
  -user string
    	Database user name
  -lockfile filepath
        Creates a lock file
  -timeout int
    	Timeout for waiting in seconds
  -timeperiod int
    	Time between checks in seconds
```
    	
The program will wait `timeperiod` seconds between each check. 'waitfordb' will finish with an exit code after the `timeout`, if no connection is available.

Exit Codes
-------------------------
| Code | Description |
|----:|------------------------|
| 10 | No database available |
|  1 | Database without tables found |
|  0 | Database with tables found |
| 101 | JDBC url is not configured |
| 103 | User is not configured |
| 104 | Password is not configured |
| 106 | Parameter timeperiod is bigger or equal to timeout. The timeout configuration must be bigger than the timeperiod. |
| 107 | Parameter timeout must be bigger than 0. |
| 108 | Parameter timeperiod must be bigger than 0. |