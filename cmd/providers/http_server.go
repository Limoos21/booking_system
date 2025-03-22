package providers

import (
	"booking_system/cmd/providers/middelware"
	"booking_system/internal/app/ports"
	"booking_system/internal/infrastructure/adapters/routers"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type HTTPServer struct {
	port        int
	Server      *gin.Engine
	logLvl      string
	controllers ports.IController
}

func NewHTTPServer(port int, logLvl string, controllers ports.IController) *HTTPServer {

	switch logLvl {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()
	corsConfig := cors.Default()
	server.Use(corsConfig)
	return &HTTPServer{
		port:        port,
		Server:      server,
		logLvl:      logLvl,
		controllers: controllers,
	}
}

func (s *HTTPServer) Run(logger *slog.Logger, jwt *middelware.Jwt) error {
	routers.New(s.Server, logger, s.controllers, jwt)
	err := s.Server.Run(fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	return nil
}
