package handler

import (
	"core/internal/models"
	"errors"
)

func (h *Handler) CheckAuth(rawToken string) (*models.UserService, error) {
	if len(rawToken) == 0 {
		return nil, errors.New("not token required")
	}

	uid, isAuth := h.Service.Account.IsAuth(rawToken)
	if !isAuth {
		return nil, errors.New("bad token required")
	}

	user := h.Service.Account.GetUserByID(int64(uid))
	if user.ID == 0 {
		return nil, errors.New("user not exist")
	}

	return &user, nil
}
