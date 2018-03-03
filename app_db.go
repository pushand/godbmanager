package godbmanager

import _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"fmt"
)

var (
	Db  *sql.DB
	err error
)

//Struct that holds sql connection details
type SQLConfig struct {
	Username string
	Password string
	Host string
	Port int
	DatabaseName string
}


//function to create connection string "root:<password>@/<dbname>"
func (sqlConfig SQLConfig) dbConnectionString() string {
	var cred string
	// [username[:password]@]
	if sqlConfig.Username != "" {
		cred = sqlConfig.Username
		if sqlConfig.Password != "" {
			cred = cred + ":" + sqlConfig.Password
		}
		cred = cred + "@"
	}
	//root:<password>@/<dbname>
	return fmt.Sprintf("%stcp([%s]:%d)/%s", cred, sqlConfig.Host, sqlConfig.Port, sqlConfig.DatabaseName)
}

//Connect to MqSql Database
func StartMySqlService(config SQLConfig) {
	Db, err = sql.Open("mysql", config.dbConnectionString())
	if err != nil {
		panic(err.Error())
	}
	//defer db.Close()
	err = Db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("DB connected")
}

//Stop MySql Database
func StoptMySqlService() {
	//if stop db is immediately called make sure than there is another instance of server running
	// find that instance with ps -ef | grep oneapp and kill
	fmt.Println("stop DB")
	Db.Close()
}
