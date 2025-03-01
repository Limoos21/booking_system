package api_server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	Server *gin.Engine
	Addr   string
}

var config = cors.Config{
	AllowOrigins:     []string{"*"},                            // Разрешенные источники
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Разрешенные методы
	AllowHeaders:     []string{"Origin", "Content-Type"},       // Разрешенные заголовки
	ExposeHeaders:    []string{"Content-Length"},               // Заголовки, доступные клиенту
	AllowCredentials: true,                                     // Разрешить передачу куки и заголовков авторизации
	MaxAge:           12 * time.Hour,                           // Время кэширования CORS
}

func NewServer(addr string) *Server {
	server := &Server{
		Server: gin.New(),
		Addr:   addr,
	}
	server.Server.Use(cors.New(config))
	return server
}
func (s *Server) SetGlobalPrefix(prefix string) *gin.RouterGroup {
	routerGroup := s.Server.Group(prefix)
	return routerGroup
}

func (s *Server) Run() error {

	return s.Server.Run()
}
