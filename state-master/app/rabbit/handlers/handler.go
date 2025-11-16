package handlers

import "git.abanppc.com/farin-project/state/app/rabbit/consumers"

type Handler interface {
	RegisterConsumer(c consumers.Consumer)
}
