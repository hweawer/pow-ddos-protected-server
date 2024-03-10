package pkg

import (
	"bufio"
	"go.uber.org/zap"
	"net"
)

type WordOfWisdomClient struct {
	address string
	logger  *zap.Logger
}

func NewWordOfWisdomClient(address string, logger *zap.Logger) *WordOfWisdomClient {
	return &WordOfWisdomClient{address, logger}
}

func (c *WordOfWisdomClient) GetWordOfWisdom() (string, error) {
	con, err := net.Dial("tcp", c.address)
	if err != nil {
		return "", err
	}
	w := bufio.NewWriter(con)
	r := bufio.NewReader(con)
	hashcash, err := c.requestChallenge(w, r)
	if err != nil {
		return "", err
	}
	if err := ComputeHashcash(hashcash); err != nil {
		return "", err
	}
	quote, err := c.requestResource(hashcash, w, r)
	if err != nil {
		return "", err
	}
	if err := c.requestStop(w); err != nil {
		return "", err
	}
	return quote, nil
}

func (c *WordOfWisdomClient) requestChallenge(w *bufio.Writer, r *bufio.Reader) (*HashcashDto, error) {
	requestChallenge := NewMessage(MesRequestChallenge, "")
	c.logger.Debug("sending request challenge", zap.Any("message", requestChallenge))
	if _, err := w.WriteString(requestChallenge.String() + "\n"); err != nil {
		return nil, err
	}
	w.Flush()
	challengeResp, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	c.logger.Debug("received challenge response", zap.String("message", challengeResp))
	challengeMes, err := ParseResponse(challengeResp)
	if err != nil {
		return nil, err
	}
	hashcash, err := Unmarshall(challengeMes.Body)
	if err != nil {
		return nil, err
	}
	return hashcash, nil
}

func (c *WordOfWisdomClient) requestResource(hashcash *HashcashDto, w *bufio.Writer, r *bufio.Reader) (string, error) {
	requestResource := NewMessage(MesRequestResource, hashcash.Marshall())
	c.logger.Debug("sending request resource", zap.Any("message", requestResource))
	if _, err := w.WriteString(requestResource.String() + "\n"); err != nil {
		return "", err
	}
	w.Flush()
	resourceResp, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	c.logger.Debug("received resource response", zap.String("message", resourceResp))
	resourceMes, err := ParseResponse(resourceResp)
	if err != nil {
		return "", err
	}
	return resourceMes.Body, nil
}

func (c *WordOfWisdomClient) requestStop(w *bufio.Writer) error {
	requestStop := NewMessage(MesStop, "")
	c.logger.Debug("sending request stop", zap.Any("message", requestStop))
	if _, err := w.WriteString(requestStop.String() + "\n"); err != nil {
		return err
	}
	w.Flush()
	return nil
}
