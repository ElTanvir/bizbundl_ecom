package admin

import (
	"bizbundl/internal/server"
)

func Init(server *server.Server) {
	router := server.GetRouter()
	adminRouter := router.Group("/admin")

	for _, route := range Routes {
		adminRouter.Add(route.Method, route.Path, route.Handler)
	}
}
