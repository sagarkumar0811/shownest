package handlers

import (
	"github.com/shownest/merchant-service/internal/usecases"
)

type Handler struct {
	usecase *usecases.UseCase
}

func New(usecase *usecases.UseCase) *Handler {
	return &Handler{usecase: usecase}
}
