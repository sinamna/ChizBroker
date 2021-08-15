package repository

import (
	"database/sql"
	"therealbroker/pkg/broker"
)

var DB PostgresDB


type PostgresDB struct {
	client *sql.DB
}
func (db *PostgresDB) SaveMessage(id int, msg broker.Message)error{
	return nil
}
func (db *PostgresDB) FetchMessage(id int)(broker.Message,error){
	return broker.Message{}, nil
}
func (db *PostgresDB) DeleteMessage(id int)error{

}

