package handler

import (
	"rip/internal/app/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}
