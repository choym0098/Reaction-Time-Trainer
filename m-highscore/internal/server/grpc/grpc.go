package grpc

import (
	"context"

	"net"

	pbhighscore "github.com/choym0098/Reaction-Time-Trainer/m-apis/m-highscore/v1"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Grpc struct {
	address string
	server  *grpc.Server
}

var HighScore = 99999999999.0

func NewServer(address string) *Grpc {
	return &Grpc{
		address: address,
	}
}

func (g *Grpc) SetHighScore(ctx context.Context, input *pbhighscore.SetHighScoreRequest) (*pbhighscore.SetHighScoreResponse, error) {
	log.Info().Msg("SetHighScore in m-highscore is called")
	HighScore = input.HighScore
	return &pbhighscore.SetHighScoreResponse{
		Set: true,
	}, nil
}

func (g *Grpc) GetHighScore(ctx context.Context, input *pbhighscore.GetHighScoreRequest) (*pbhighscore.GetHighScoreResponse, error) {
	log.Info().Msg("GetHighScore in m-highscore is called")
	return &pbhighscore.GetHighScoreResponse{
		HighScore: HighScore,
	}, nil
}

func (g *Grpc) ListenAndServe() error {
	listener, err := net.Listen("tcp", g.address)
	if err != nil {
		return errors.Wrap(err, "failed to open tcp port")
	}

	serverOptions := []grpc.ServerOption{}

	g.server = grpc.NewServer(serverOptions...)

	pbhighscore.RegisterGameServer(g.server, g)

	log.Info().Str("address", g.address).Msg("starting gRPC server for highscore microservice")

	err = g.server.Serve(listener)
	if err != nil {
		return errors.Wrap(err, "failed to start gRPC server for highscore microservice")
	}

	return nil
}
