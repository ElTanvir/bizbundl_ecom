package views

import (
	"bizbundl/internal/server"
	"bizbundl/internal/views/frontend"
)

// Initializes all views
// Admin views Stays Same for All Site
// Frontend views/ Customer Facing Views may Change For Each Site
// We Just Replace the Frontend views/ Customer Facing Views for Each Site If Need
// While Maintaining the Same Structure
func Init(server *server.Server) {
	frontend.Init(server)
}
