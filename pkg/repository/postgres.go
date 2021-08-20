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

func (db *PostgresDatabase) SaveMessage(id int, msg broker.Message, subject string) error {
	//rows,err := db.client.Query("call save_message($1,$2,$3,$4);", id, msg.Body, subject, int32(msg.Expiration))
	query := fmt.Sprintf(`INSERT INTO messages(id, subject, body, expiration_date)VALUES (%d, '%s', '%s', %v);`,
		id, msg.Body, subject, int32(msg.Expiration))
	//fmt.Println(query)
	rows, err := db.client.Query(query)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	rows.Close()
	fmt.Println("saved")

	return nil
}
func (db *PostgresDatabase) FetchMessage(id int, subject string) (broker.Message, error) {
	query := fmt.Sprintf("SELECT body, expiration_date from messages where messages.id=%d and messages.subject=%s;",
		id, subject)
	rows, err := db.client.Query(query)

	if err != nil {
		return broker.Message{}, err
	}
	var body string
	var expirationDate int32
	err = rows.Scan(&body, &expirationDate)
	if err != nil {
		return broker.Message{}, err
	}
	msg := broker.Message{
		Body:       body,
		Expiration: time.Duration(expirationDate),
	}
	rows.Close()
	return msg, nil

}
func (db *PostgresDatabase) DeleteMessage(id int, subject string) error {
	query := fmt.Sprintf(`DELETE FROM messages WHERE messages.id=%d and messages.subject='%s';`, id, subject)
	rows, err := db.client.Query(query)

	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func GetPostgreDB() (Database, error) {
	var once sync.Once
	once.Do(func() {
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB"))
		//fmt.Println(connString)
		client, err := sql.Open("postgres", connString)
		if err != nil {
			connectionError = err
			return
		}
		//defer client.Close()
		err = client.Ping()
		if err != nil {
			connectionError = err
			return
		}
		fmt.Println("connected to postgres.")
		err = createTable(client)
		if err != nil {
			connectionError = err
			return
		}
		client.SetMaxOpenConns(90)
		postgresDB = &PostgresDatabase{client: client}
	})
	return postgresDB, connectionError
}

func createTable(client *sql.DB) error {
	table := `
	CREATE TABLE IF NOT EXISTS messages (
		id integer not null,
		subject varchar(255) not null,
		body varchar(255) not null,
		expiration_date bigint not null,
		primary key(id)
);
`
	_, err := client.Exec(table)
	if err != nil {
		return err
	}
	return nil
}
