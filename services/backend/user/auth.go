package user

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Login(ctx context.Context, rc *redis.Client, ui int64, exp time.Duration) (string, error) {
	n := time.Now()
	rsi := fmt.Sprintf("%d%d", ui, n.UnixNano())
	s := sha512.Sum512([]byte(rsi))
	si := base64.StdEncoding.EncodeToString(s[:])

	if err := rc.Set(ctx, si, fmt.Sprintf("%d", ui), exp).Err(); err != nil {
		return "", err
	}
	return si, nil
}

func Logout(ctx context.Context, rc *redis.Client, si string) error {
	return rc.Del(ctx, si).Err()
}
