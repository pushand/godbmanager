package oneappconfig

import _ "github.com/go-sql-driver/mysql"
import (
	"database/sql"
	"fmt"
)

var (
	Db  *sql.DB
	err error
	sqlConfig SQLConfig
)

type SQLConfig struct {
	Username string
	Password string
	Host string
	Port int
	DatabaseName string
}

func init() {
	Db, err = sql.Open("mysql", sqlConfig.dbConnectionString())
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

func Start(config SQLConfig) {
	sqlConfig = config
}

func Stop() {
	//if stop db is immediately called make sure than there is another instance of server running
	// find that instance with ps -ef | grep oneapp and kill
	fmt.Println("stop DB")
	Db.Close()
}
