package models



type DBMAP map[string]interface{}


var Mysqlconn=DBMAP{

	"username": "root",
	"password": "root",
	"host":     "localhost",
	"port":     3306,
	"name":     "test",
}

var Orclconn=DBMAP{
	"fmt": "null",
	"username": "system",
	"password": "123456",
	"host":     "127.0.0.1",
	"name":     "orcl",
}