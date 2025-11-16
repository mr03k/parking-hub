package handlers

import "farin/app/rabbit/consumers"

type Handler interface {
	RegisterConsumer(c consumers.Consumer)
}
