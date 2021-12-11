package user_session

import (
	"github.com/google/uuid"
	"main/pkg/gredis"
	"strconv"
	"time"
)

// auth cache namespace
const authCachePrefix = "auth::"

// OpenSession establish a link between sessionID and userID
func OpenSession(userID uint) (string, error) {
	sessionID := uuid.New().String()
	if err := gredis.Set(authCachePrefix+sessionID, strconv.Itoa(int(userID)), 30*time.Minute); err != nil {
		return "", err
	}
	return sessionID, nil
}

// CloseSession remove the link between sessionID and userID
func CloseSession(sessionID string) error {
	if err := gredis.Delete(authCachePrefix + sessionID); err != nil {
		return err
	}
	return nil
}

func GetUserID(sessionID string) (uint, error) {
	v, err := gredis.Get(authCachePrefix + sessionID)
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}
