Wait For DB (MS SQL)
=========================
This tool checks the configured database connection. If the database is not available the command exits with a special exit code.

Usage
-------------------------

It is possible to configure the following parameters:

Usage of waitfordb:
  -host string
    	Host name of the server with a MS SQLServer (default "localhost")
  -name string
    	Database name
  -password string
    	Passwort of database user
  -port int
    	MSSQLServer is listing on this port (default 1433)
  -user string
    	Database user name
    	
The program will wait 20 seconds between each check. It tries 10 times to connect the database. 
This is currently not configurable.

Exit Codes
-------------------------
| Code | Description |
|----:|------------------------|
| 10 | No database available |
|  1 | Database without tables found |
|  0 | Database with tables found |
| 101 | Host is not configured |
| 102 | Port is not configured |
| 103 | User is not configured |
| 104 | Password is not configured |
| 105 | Database name is not configured |