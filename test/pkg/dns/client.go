package dns

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
)

type Client struct {
	client dns.Client
}

func NewClient() *Client {
	return new(Client)
}

func (c *Client) QueryTypeA(serverIP, serverPort, name string) (
	answer []string, e error,
) {
	const (
		addressFormat = "%s:%s"
	)

	var (
		address  string
		answerIP net.IPAddr
		answerRR dns.RR
		message  *dns.Msg
	)

	address = fmt.Sprintf(addressFormat, serverIP, serverPort)

	message = new(dns.Msg)

	message.SetQuestion(name, dns.TypeA)

	message, _, e = c.client.Exchange(message, address)
	if e != nil {
		return
	}

	for _, answerRR = range message.Answer {
		answerIP = net.IPAddr{
			IP: answerRR.(*dns.A).A,
		}

		answer = append(answer,
			answerIP.String(),
		)
	}

	return
}
