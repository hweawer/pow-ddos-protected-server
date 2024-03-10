package pkg

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrChallengeExpired = "challenge expired"
	ErrInvalidResource  = "invalid resource"
	ErrInvalidHash      = "invalid hash"
)

// HashcashDto represents a hashcash challenge.
type HashcashDto struct {
	Version    int
	ZerosCount int
	Date       int64
	// Client IP
	Resource string
	Rand     string
	Counter  int
}

// NewHashcashDto creates a new instance of HashcashDto.
func NewHashcashDto(resource string, rand string) *HashcashDto {
	return &HashcashDto{
		Version: 1,
		// todo: can pass from config
		ZerosCount: 3,
		Date:       time.Now().Unix(),
		Resource:   resource,
		Rand:       rand,
		Counter:    0,
	}
}

// String returns the string representation of the hashcash challenge.
func (h HashcashDto) String() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.ZerosCount, h.Date, h.Resource, h.Rand, h.Counter)
}

// IsValid checks if the hashcash challenge is valid.
func (h HashcashDto) IsValid(resource string) error {
	if time.Unix(h.Date, 0).Add(time.Duration(2)*time.Hour*24).Unix() < time.Now().Unix() {
		return fmt.Errorf(ErrChallengeExpired)
	}
	if resource != h.Resource {
		return fmt.Errorf(ErrInvalidResource)
	}

	if !h.isHashCorrect() {
		return fmt.Errorf(ErrInvalidHash)
	}

	return nil
}

func (h HashcashDto) sha1Hash() string {
	return sha1Hash(h.String())
}

func (h HashcashDto) isHashCorrect() bool {
	hash := h.sha1Hash()
	zeros := h.ZerosCount
	if zeros > len(hash) {
		return false
	}
	for _, ch := range hash[:zeros] {
		if string(ch) != "0" {
			return false
		}
	}
	return true
}

// ComputeHashcash computes the hashcash challenge.
func ComputeHashcash(hashCashDto *HashcashDto) error {
	// todo: can pass from config
	limit := 10000000
	for i := 0; i < limit; i++ {
		if hashCashDto.isHashCorrect() {
			return nil
		}
		hashCashDto.Counter++
	}
	return fmt.Errorf("max iterations exceeded")
}

func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (h HashcashDto) Marshall() string {
	r := base64.StdEncoding.EncodeToString([]byte(h.Rand))
	c := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(h.Counter)))
	// it is an IP, need to encode to parse it properly
	a := base64.StdEncoding.EncodeToString([]byte(h.Resource))
	s := fmt.Sprintf("%d:%d:%d:%s::%s:%s", h.Version, h.ZerosCount, h.Date, a, r, c)
	return s
}

func Unmarshall(s string) (*HashcashDto, error) {
	s = strings.Trim(s, "\n")
	s = strings.TrimSpace(s)
	tokens := strings.Split(s, ":")
	if len(tokens) != 7 {
		return nil, fmt.Errorf("invalid hashcash")
	}
	var hashcash HashcashDto
	var err error
	hashcash.Version, err = strconv.Atoi(tokens[0])
	if err != nil {
		return nil, err
	}
	hashcash.ZerosCount, err = strconv.Atoi(tokens[1])
	if err != nil {
		return nil, err
	}
	hashcash.Date, err = strconv.ParseInt(tokens[2], 10, 64)
	if err != nil {
		return nil, err
	}
	ba, err := base64.StdEncoding.DecodeString(tokens[3])
	if err != nil {
		return nil, err
	}
	hashcash.Resource = string(ba)
	br, err := base64.StdEncoding.DecodeString(tokens[5])
	if err != nil {
		return nil, err
	}
	hashcash.Rand = string(br)

	bc, err := base64.StdEncoding.DecodeString(tokens[6])
	if err != nil {
		return nil, err
	}
	hashcash.Counter, err = strconv.Atoi(string(bc))
	if err != nil {
		return nil, err
	}
	return &hashcash, nil
}
