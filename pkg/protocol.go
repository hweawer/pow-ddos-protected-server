package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MesStop = iota
	MesRequestChallenge
	MesResponseChallenge
	MesRequestResource
	MesResponseResource
)

const Delimiter = "|"

type Message struct {
	Type int
	Body string
}

func (m *Message) String() string {
	return fmt.Sprintf("%d%s%s", m.Type, Delimiter, m.Body)
}

func NewMessage(t int, b string) *Message {
	return &Message{Type: t, Body: b}
}

func parseMessage(s string) (*Message, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, Delimiter)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid message")
	}
	t, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid message type")
	}
	// MesStop and MesRequestChallenge don't have a body
	if t == MesStop {
		return &Message{Type: MesStop}, nil
	}
	if t == MesRequestChallenge {
		return &Message{Type: MesRequestChallenge}, nil
	}
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid message")
	}
	return &Message{Type: t, Body: parts[1]}, nil
}

func ParseRequest(s string) (*Message, error) {
	mes, err := parseMessage(s)
	if err != nil {
		return nil, err
	}
	if mes.Type != MesRequestResource && mes.Type != MesRequestChallenge && mes.Type != MesStop {
		return nil, fmt.Errorf("invalid message type")
	}
	return mes, nil
}

func ParseResponse(s string) (*Message, error) {
	mes, err := parseMessage(s)
	if err != nil {
		return nil, err
	}
	if mes == nil {
		return nil, fmt.Errorf("invalid message")
	}
	if mes.Type != MesResponseResource && mes.Type != MesResponseChallenge {
		return nil, fmt.Errorf("invalid message type")
	}
	return mes, nil
}
