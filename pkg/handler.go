package pkg

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"math/rand"
	"net"
	"time"
)

// todo: can pass from config
const DefaultExpiration = time.Minute

var quotes = []string{
	"Guard well your thoughts when alone and your words when accompanied.",
	"Change in all things is sweet.",
	"Time is a created thing. To say 'I don't have time,' is like saying, 'I don't want to.",
	"Be kind, for everyone you meet is fighting a hard battle.",
	"Courage is knowing what not to fear.",
	"Nothing endures but change.",
	"Time is the wisest counselor of all.",
}

var (
	ErrInvalidMessageType = "invalid message type"
)

type Handler struct {
	cache *cache.Cache
}

// NewHandler creates a new instance of Handler.
func NewHandler() *Handler {
	return &Handler{
		cache: cache.New(DefaultExpiration, 5*time.Minute),
	}
}

// HandleMessage handles the incoming message and returns the response message.
// If the message is invalid, it returns an error.
func (h *Handler) HandleMessage(c net.Conn, mes *Message) (*Message, error) {
	switch mes.Type {
	case MesRequestResource:
		return h.handleRequestResource(c, mes)
	case MesRequestChallenge:
		return h.handleRequestChallenge(c)
	default:
		return nil, fmt.Errorf(ErrInvalidMessageType)
	}
}

// handleRequestChallenge returns a new challenge for the client.
func (h *Handler) handleRequestChallenge(c net.Conn) (*Message, error) {
	r := randomString(8)
	if err := h.cache.Add(r, struct{}{}, DefaultExpiration); err != nil {
		return nil, err
	}
	hashCash := NewHashcashDto(c.RemoteAddr().String(), r)
	msg := NewMessage(MesResponseChallenge, hashCash.Marshall())
	return msg, nil
}

// handleRequestResource returns a resource for the client.
func (h *Handler) handleRequestResource(c net.Conn, mes *Message) (*Message, error) {
	hashcash, err := Unmarshall(mes.Body)
	if err != nil {
		return nil, err
	}

	clientIp := c.RemoteAddr().String()
	if err := hashcash.IsValid(clientIp); err != nil {
		return nil, err
	}

	_, found := h.cache.Get(hashcash.Rand)
	if !found {
		return nil, fmt.Errorf("didn't find challenge")
	}

	msg := NewMessage(MesResponseResource, quotes[rand.Intn(len(quotes))])
	return msg, nil
}

// randomString returns a random string of length n for Hashcash algorithm.
func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
