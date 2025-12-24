package frontend

import (
	"bizbundl/internal/server"
)

func Init(server *server.Server) {
	router := server.GetRouter()

	for _, route := range Routes {
		router.Add(route.Method, route.Path, route.Handler)
	}
}
