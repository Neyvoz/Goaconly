module goaconly

go 1.26.2

require github.com/lib/pq v1.12.3

require github.com/google/uuid v1.6.0

require golang.org/x/crypto v0.54.0

require (
	github.com/caarlos0/env/v11 v11.4.1
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/redis/go-redis/v9 v9.21.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)
