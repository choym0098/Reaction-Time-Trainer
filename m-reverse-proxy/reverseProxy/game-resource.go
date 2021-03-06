package reverseProxy

import (
	"context"
	"strconv"

	pbgameengine "github.com/choym0098/Reaction-Time-Trainer/m-apis/m-game-engine/v1"
	pbhighscore "github.com/choym0098/Reaction-Time-Trainer/m-apis/m-highscore/v1"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type gameResource struct {
	gameClient       pbhighscore.GameClient
	gameEngineClient pbgameengine.GameClient
}

func NewGameResource(gameClient pbhighscore.GameClient, gameEngineClient pbgameengine.GameClient) *gameResource {
	return &gameResource{
		gameClient:       gameClient,
		gameEngineClient: gameEngineClient,
	}
}

func NewGrpcGameServiceClient(serverAddr string) (pbhighscore.GameClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-highscore connection is nil in m-reverse-proxy")
	}

	client := pbhighscore.NewGameClient(conn)

	return client, nil
}

func NewGrpcGameEngineServiceClient(serverAddr string) (pbgameengine.GameClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("Failed to dial: %v", err)
		return nil, err
	} else {
		log.Info().Msgf("Successfully connected to [%s]", serverAddr)
	}

	if conn == nil {
		log.Info().Msg("m-game-engine connection is nil in m-reverse-proxy")
	}

	client := pbgameengine.NewGameClient(conn)

	return client, nil
}

func (gr *gameResource) SetHighScore(c *gin.Context) {
	highScoreString := c.Param("hs")
	highScoreFloat64, err := strconv.ParseFloat(highScoreString, 64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert high score string to float64")
	}

	gr.gameClient.SetHighScore(context.Background(), &pbhighscore.SetHighScoreRequest{
		HighScore: highScoreFloat64,
	})
}

func (gr *gameResource) GetHighScore(c *gin.Context) {
	highScoreResponse, err := gr.gameClient.GetHighScore(context.Background(), &pbhighscore.GetHighScoreRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting highscore")
	}

	hsString := strconv.FormatFloat(highScoreResponse.HighScore, 'e', -1, 64)

	c.JSONP(200, gin.H{
		"hs": hsString,
	})
}

func (gr *gameResource) GetSize(c *gin.Context) {
	sizeResponse, err := gr.gameEngineClient.GetSize(context.Background(), &pbgameengine.GetSizeRequest{})
	if err != nil {
		log.Error().Err(err).Msg("Error while getting size")
	}

	c.JSON(200, gin.H{
		"size": sizeResponse.GetSize(),
	})
}

func (gr *gameResource) SetScore(c *gin.Context) {
	scoreString := c.Param("score")
	scoreFloat64, _ := strconv.ParseFloat(scoreString, 64)

	_, err := gr.gameEngineClient.SetScore(context.Background(), &pbgameengine.SetScoreRequest{
		Score: scoreFloat64,
	})
	if err != nil {
		log.Error().Err(err).Msg("Error while setting score in m-game-engine")
	}
}
