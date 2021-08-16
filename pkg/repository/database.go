package repository

import "therealbroker/pkg/broker"

type Database interface{
	SaveMessage(id int, msg broker.Message, subject string)
	FetchMessage(id int, subject string) (broker.Message, error, )
	DeleteMessage(id int, subject string)
}
