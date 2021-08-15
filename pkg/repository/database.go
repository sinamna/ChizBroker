package repository

import "therealbroker/pkg/broker"

type Database interface{
	SaveMessage(id int, msg broker.Message)
	FetchMessage(id int) (broker.Message, error)
	DeleteMessage(id int)
}
