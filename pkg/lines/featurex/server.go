package featurex

import (
	"fmt"

	"github.startlite.cn/itapp/startlite/pkg/lines/constantx"
	"github.startlite.cn/itapp/startlite/pkg/lines/runlevelx"
	"github.startlite.cn/itapp/startlite/pkg/lines/typesx"
)

type Server struct {
	typesx.ServerConfig
	runlevelx.RunLevel
}

func MustProvideServer(runLevel string) *Server {

	server := &Server{
		ServerConfig: typesx.ServerConfig{
			Name:     "starlite",
			Port:     8080,
			RunLevel: runLevel,
			Host:     "127.0.0.1",
		},
		RunLevel: runlevelx.RunLevelLocal,
	}
	//cl.Load(server)
	server.RunLevel = runlevelx.DetermineRunLevel(server.ServerConfig.RunLevel)
	return server
}

// WhoAmI introduces application itself
func (server *Server) WhoAmI() string {
	return server.ServerConfig.Name
}

func (server Server) Host() string {
	baseURL := server.ServerConfig.Host
	if len(baseURL) == 0 {
		baseURL = "http://localhost:4200"
	}
	return baseURL
}

func (server *Server) APIEndpoint() string {
	return fmt.Sprintf("%s%s", server.Host(), constantx.APIPrefix)
}
