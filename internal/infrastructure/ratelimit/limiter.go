package ratelimit

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// incrementScript атомарно увеличивает счётчик и выставляет TTL только
// при первом инкременте в окне. Это критично: если делать INCR и EXPIRE
// двумя отдельными командами, между ними возможна гонка — если процесс
// упадёт между INCR и EXPIRE, ключ останется висеть без TTL навсегда.
var incrementScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
    redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return current
`)

type RedisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{client: client}
}

// Allow проверяет, не превышен ли лимит запросов для данного ключа
// в пределах окна windowSeconds. Возвращает (разрешено, текущий счётчик, ошибка)
func (l *RedisLimiter) Allow(ctx context.Context, key string, limit int, windowSeconds int) (bool, int, error) {
	result, err := incrementScript.Run(ctx, l.client, []string{key}, windowSeconds).Int()
	if err != nil {
		return false, 0, fmt.Errorf("ratelimit: script exec failed: %w", err)
	}

	return result <= limit, result, nil
}
