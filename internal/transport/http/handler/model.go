package handler

import (
	"core/internal/services"
	"core/internal/tg/bot"
)

type Handler struct {
	Service *services.Services
	Bot     *bot.Bot
}
