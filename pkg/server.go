package pkg

import (
	"bufio"
	"go.uber.org/zap"
	"net"
	"strings"
)

type Server struct {
	handler  *Handler
	listener net.Listener
	logger   *zap.Logger
}

// NewServer creates a new server
func NewServer(handler *Handler, listener net.Listener, logger *zap.Logger) *Server {
	return &Server{handler, listener, logger}
}

// Serve starts the server
func (s *Server) Serve() {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	s.logger.Debug("handling connection", zap.String("remote", c.RemoteAddr().String()))
	w := bufio.NewWriter(c)
	r := bufio.NewReader(c)
	for {
		netData, err := r.ReadString('\n')
		if err != nil {
			s.logger.Error("failed to read message", zap.Error(err))
			continue
		}
		s.logger.Debug("received message", zap.String("message", netData))

		temp := strings.TrimSpace(netData)
		message, err := ParseRequest(temp)
		if err != nil {
			s.logger.Error("failed to parse message", zap.Error(err))
			continue
		}
		s.logger.Debug("parsed message", zap.Any("message", message))
		if message == nil {
			if _, err := w.WriteString("unknown message\n"); err != nil {
				s.logger.Error("failed to write message", zap.Error(err))
			}
			w.Flush()
			continue
		}
		if message.Type == MesStop {
			break
		}
		resp, err := s.handler.HandleMessage(c, message)
		if resp != nil {
			s.logger.Debug("sending response", zap.Any("response", resp))
			if _, err := w.WriteString(resp.String() + "\n"); err != nil {
				s.logger.Error("failed to write message", zap.Error(err))
			}
		}
		if err != nil {
			s.logger.Debug("sending response", zap.Any("response", err))
			if _, err := w.WriteString(err.Error() + "\n"); err != nil {
				s.logger.Error("failed to write message", zap.Error(err))
			}
		}
		w.Flush()
	}
	c.Close()
}
