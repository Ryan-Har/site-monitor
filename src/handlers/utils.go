package handlers

import (
	"context"
	"errors"
	"github.com/Ryan-Har/site-monitor/src/models"
)

func GetUserInfoFromContext(ctx context.Context) (models.UserInfo, error) {
	userInfo, ok := ctx.Value(models.UserInfoKey).(models.UserInfo)
	if !ok {
		return models.UserInfo{}, errors.New("no user info found in context")
	}

	return userInfo, nil
}
