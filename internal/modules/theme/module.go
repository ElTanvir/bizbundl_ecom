package theme

import (
	"context"
	"portfolioed/internal/modules/theme/repository"
	"portfolioed/internal/modules/theme/usecase"
	"portfolioed/internal/server"
)


func Init(server *server.Server) {
	useCase := getUseCase(server)
	useCase.InvalidateCache(context.Background())
}

func getUseCase(s *server.Server) usecase.ThemeUseCase {
	repo := repository.NewThemeRepository(s.GetDB())
	useCase := usecase.NewThemeUseCase(repo)
	return useCase
}

func GetHandler(s *server.Server) usecase.ThemeUseCase {
	useCase := getUseCase(s)
	return useCase
}

