package handlers

import "github.com/shownest/pricing-service/internal/usecases"

type Handler struct {
	usecase *usecases.UseCase
}

func New(uc *usecases.UseCase) *Handler {
	return &Handler{usecase: uc}
}
