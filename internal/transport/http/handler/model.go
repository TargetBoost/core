package handler

import (
	"core/internal/services"
	"core/internal/transport/tg/bot"
)

type Handler struct {
	Service *services.Services
	Bot     *bot.Bot
}
