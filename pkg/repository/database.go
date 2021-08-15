package repository

import "therealbroker/pkg/broker"

type Database interface{
	SaveMessage(id int, msg broker.Message)error
	FetchMessage(id int) (broker.Message, error)
	DeleteMessage(id int) error
}
