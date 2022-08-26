package transport

import (
	"context"
	"net"
	"time"
)

type Dialer struct{}

func NewDialer() *Dialer {
	return new(Dialer)
}

func (d *Dialer) Dial(network, address string, timeout time.Duration) (
	e error,
) {
	var (
		dialer net.Dialer
		timer  context.Context
	)

	timer, _ = context.WithTimeout(
		context.Background(),
		timeout,
	)

	for {
		_, e = dialer.Dial(network, address)
		if e == nil {
			return
		}

		select {
		case <-timer.Done():
			return

		default:
			break
		}
	}

	return
}
