package repository

import "therealbroker/pkg/broker"

type Database interface{
	SaveMessage(id int, msg broker.Message, subject string)error
	FetchMessage(id int, subject string) (broker.Message, error)
	DeleteMessage(id int, subject string)error
}
