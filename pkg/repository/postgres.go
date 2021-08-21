package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	//"strconv"

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

func (db *PostgresDatabase) SaveMessage (msg broker.Message, subject string) int{
	//fmt.Println("msg:",msg)
	var insertedID int
	err := db.client.QueryRow(`INSERT INTO messages(id, subject, body, expiration_date) VALUES (DEFAULT, $1, $2, $3) RETURNING id;`,subject, msg.Body, int32(msg.Expiration)).Scan(&insertedID)
	if err != nil {
		fmt.Println("saving error:",err)
		return -1
	}
	//fmt.Println("saved")
	return insertedID
}
func (db *PostgresDatabase) FetchMessage(id int, subject string) (broker.Message, error) {
	query := fmt.Sprintf("SELECT body, expiration_date from messages where messages.id=%d and messages.subject='%s';",
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
func (db *PostgresDatabase) DeleteMessage(id int, subject string) {
	query := fmt.Sprintf(`DELETE FROM messages WHERE messages.id=%d and messages.subject='%s';`, id, subject)
	_, err := db.client.Exec(query)

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("deleted")
	//rows.Close()
	//return nil
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
		err = createIndex(client)
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
		id serial,
		subject varchar(255) not null,
		body varchar(255) ,
		expiration_date bigint not null,
		primary key(id, subject)
);
`
	_, err := client.Exec(table)
	if err != nil {
		return err
	}
	return nil
}
func createIndex(client *sql.DB) error {
	command := `CREATE INDEX IF NOT EXISTS idx_id_subject on messages (id,subject)`
	_, err := client.Exec(command)
	if err != nil {
		return err
	}
	return nil
}
