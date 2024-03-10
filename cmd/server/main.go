package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"pow_server/pkg"
)

// NewListener creates a new net.Listener.
func NewListener(lc fx.Lifecycle) (net.Listener, error) {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})
	return l, nil
}

func main() {
	fx.New(
		fx.Provide(NewListener, pkg.NewHandler, pkg.NewServer, zap.NewExample),
		fx.Invoke(func(s *pkg.Server) { s.Serve() }),
	).Run()
}
