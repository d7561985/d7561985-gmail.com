//go:generate mockery  -all
package delivery

import (
	"net"
)

type Interface interface {
	Serve(listener net.Listener)
	Stop()
}
