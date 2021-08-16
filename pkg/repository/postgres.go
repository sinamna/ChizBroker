package repository

import (
	"database/sql"
	"fmt"
	"github.com/prometheus/common/log"
	"os"
	"sync"
	"therealbroker/pkg/broker"
	"time"
)

var postgresDB *PostgresDatabase
var connectionError error

type PostgresDatabase struct {
	client *sql.DB
}
func (db *PostgresDatabase) SaveMessage(id int, msg broker.Message, subject string){
	_,err := db.client.Query("call save_message(?,?,?,?)", id, msg.Body, subject, int32(msg.Expiration))
	if err!= nil{
		log.Errorln(err)
	}
}
func (db *PostgresDatabase) FetchMessage(id int, subject string)(broker.Message,error){
	query := fmt.Sprintf("SELECT body, expiration_date from messages where messages.id=%d and messages.subject=%s",
		id,subject)
	rows, err := db.client.Query(query)
	if err!= nil{
		return broker.Message{}, err
	}else{
		var body string
		var expirationDate int32
		err := rows.Scan(&body, &expirationDate)
		if err != nil {
			return broker.Message{}, err
		}
		msg := broker.Message{
			Body: body,
			Expiration: time.Duration(expirationDate),
		}
		return msg, nil
	}
}
func (db *PostgresDatabase) DeleteMessage(id int,subject string){

}

func GetPostgreDB()(Database,error){
	var once sync.Once
	once.Do(func() {
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("host"),os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
		//fmt.Println(connString)
		client, err:= sql.Open("postgres",connString)
		if err!=nil{
			connectionError=err
			return
		}
		defer client.Close()
		err = client.Ping()
		if err != nil {
			connectionError=err
			return
		}
		postgresDB = &PostgresDatabase{client: client}
	})
	return postgresDB,connectionError
}

