package handler

import (
	"rip/internal/app/redis"
	"rip/internal/app/repository"
)

type Handler struct {
	Repo        *repository.Repository
	RedisClient *redis.Client
}

func New(repo *repository.Repository, redisClient ...*redis.Client) *Handler {
	h := &Handler{Repo: repo}
	if len(redisClient) > 0 {
		h.RedisClient = redisClient[0]
	}
	return h
}
