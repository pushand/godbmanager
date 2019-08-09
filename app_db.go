package godbmanager

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
)
import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	Db  *sql.DB
	err error
)

//Struct that holds sql connection details
type SQLConfig struct {
	Username     string
	Password     string
	Host         string
	Port         int
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

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

//Connect to MqSql Database
func StartMySqlService(config SQLConfig) {

	if config.Host != "" {
		Db, err = sql.Open("mysql", config.dbConnectionString())
	} else {
		//user:password@tcp(IP:3306)/dbname
		fmt.Println("CLOUD_SQL", os.Getenv("CLOUD_SQL"))
		Db, err = sql.Open("mysql", os.Getenv("CLOUD_SQL"))
	}

	if err != nil {
		log.Panic(err.Error())
	}
	//defer db.Close()
	err = Db.Ping()
	if err != nil {
		log.Panic(err.Error())
	}
	fmt.Println("DB connected")
}

func StartMySqlSSHService(sshHost, sshKey, sshUser, sshKeyPass, dbHostPort, dbName, dbUser, dbPass string) {

	//var hostKey ssh.PublicKey

	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(sshKey, sshKeyPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if sshcon, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshHost, 22), sshConfig); err == nil {
		//defer sshcon.Close()

		fmt.Println("ssh success connected")
		// Now we register the ViaSSHDialer with the ssh connection as a parameter
		mysql.RegisterDial("mysql+tcp", (&ViaSSHDialer{sshcon}).Dial)

		Db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", dbUser, dbPass, dbHostPort, dbName))

		if err != nil {
			log.Panic(err.Error())
		}
		//defer db.Close()
		err = Db.Ping()
		if err != nil {
			log.Panic(err.Error())
		}
		fmt.Println("DB connected")
	} else {
		log.Panic(err.Error())
	}
}

func PublicKeyFile(file, pass string) ssh.AuthMethod {

	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.Panic(err.Error())
		return nil
	}

	fmt.Println("key pass", pass)
	key, err := ssh.ParsePrivateKey(decrypt(buffer, []byte(pass)))
	if err != nil {
		log.Panic(err.Error())
		return nil
	}
	return ssh.PublicKeys(key)
}

func decrypt(key []byte, password []byte) []byte {
	block, rest := pem.Decode(key)
	if len(rest) > 0 {
		log.Fatalf("Extra data included in key")
	}

	if x509.IsEncryptedPEMBlock(block) {
		der, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			log.Fatalf("Decrypt failed: %v", err)
		}
		return pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der})
	}
	return key
}

//Stop MySql Database
func StoptMySqlService() {
	//if stop db is immediately called make sure than there is another instance of server running
	// find that instance with ps -ef | grep oneapp and kill
	fmt.Println("stop DB")
	Db.Close()
}
