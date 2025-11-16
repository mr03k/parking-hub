package handlers

import "git.abanppc.com/farin-project/vehicle-records/app/rabbit/consumers"

type Handler interface {
	RegisterConsumer(c consumers.Consumer)
}
