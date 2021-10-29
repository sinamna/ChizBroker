package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"

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
	sync.Mutex
	client         *sql.DB
	addMessages    []string
	deleteMessages []string
}

func (db *PostgresDatabase) SaveMessage(msg broker.Message, subject string) int {
	query := fmt.Sprintf(`INSERT INTO messages(id, subject, body, expiration_date) VALUES (DEFAULT, '%s', '%s', %v) RETURNING id;`,subject, msg.Body, int64(msg.Expiration))
	var insertedID int
	row, err := db.client.Query(query)
	row.Next()
	row.Scan(&insertedID)
	if err != nil {
		fmt.Println("saving error:", err)
		return -1
	}
	row.Close()
	return insertedID
}
func (db *PostgresDatabase) FetchMessage(id int, subject string) (broker.Message, error) {
	query := fmt.Sprintf("SELECT body, expiration_date from messages where messages.id=%d and messages.subject='%s';",
		id, subject)
	rows, err := db.client.Query(query)

	if err != nil {
		fmt.Println("fetch: returned from query")
		return broker.Message{}, err
	}
	var body string
	var expirationDate int64
	for rows.Next() {
		err = rows.Scan(&body, &expirationDate)
		if err != nil {
			fmt.Println("fetch: scan error")
			return broker.Message{}, err
		}
	}
	if err := rows.Err();err!=nil{
		fmt.Println("rows err: ",err)
	}
	msg := broker.Message{
		Body:       body,
		Expiration: time.Duration(expirationDate),
	}
	//rows.Close()
	return msg, nil

}
func (db *PostgresDatabase) DeleteMessage(id int, subject string) {
	db.Lock()
	db.deleteMessages = append(db.deleteMessages,fmt.Sprintf("(id,subject)=(%d,'%s')",id,subject))
	db.Unlock()
}
func (db *PostgresDatabase) batchOperationHandler(ticker *time.Ticker){
	for {
		select {
		case <- ticker.C:
			db.Lock()
			if len(db.deleteMessages) != 0{
				query := `DELETE FROM public.messages WHERE ` + strings.Join(db.deleteMessages," or ")+";"
				db.deleteMessages = db.deleteMessages[:0]
				_,err := db.client.Exec(query)
				if err!= nil{
					fmt.Println(err)
				}
			}
			db.Unlock()
		}
	}
}
func GetPostgreDB() (Database, error) {
	var once sync.Once
	once.Do(func() {
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("PUSER"), os.Getenv("PASSWORD"), os.Getenv("DB"))
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
		client.SetMaxIdleConns(45)
		client.SetConnMaxIdleTime(time.Second*10)
		postgresDB = &PostgresDatabase{
			client:         client,
			addMessages:    make([]string, 0),
			deleteMessages: make([]string, 0),
		}
		ticker := time.NewTicker(1 * time.Second)
		go postgresDB.batchOperationHandler(ticker)
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
