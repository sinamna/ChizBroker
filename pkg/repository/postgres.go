package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	//"github.com/prometheus/common/log"
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
func (db *PostgresDatabase) SaveMessage(id int, msg broker.Message, subject string)error{
	fmt.Println("saving")
	rows,err := db.client.Query("call save_message($1,$2,$3,$4);", id, msg.Body, subject, int32(msg.Expiration))
	defer rows.Close()
	if err!= nil{
		return err
	}
	return nil
}
func (db *PostgresDatabase) FetchMessage(id int, subject string)(broker.Message,error){
	query := fmt.Sprintf("SELECT body, expiration_date from messages where messages.id=%d and messages.subject=%s;",
		id,subject)
	rows, err := db.client.Query(query)
	defer rows.Close()
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
func (db *PostgresDatabase) DeleteMessage(id int,subject string)error{
	rows, err:= db.client.Query("call delete_message($1,$2);",id,subject)
	defer rows.Close()
	if err!=nil{
		return err
	}
	return nil
}

func GetPostgreDB()(Database,error){
	var once sync.Once
	once.Do(func() {
		fmt.Println("new DB connection")
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("host"),os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
		//fmt.Println(connString)
		client, err:= sql.Open("postgres",connString)
		if err!=nil{
			connectionError=err
			return
		}
		//defer client.Close()
		err = client.Ping()
		if err != nil {
			connectionError=err
			return
		}
		fmt.Println("connected to postgres.")
		postgresDB = &PostgresDatabase{client: client}
	})
	return postgresDB,connectionError
}

