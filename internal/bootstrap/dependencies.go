package bootstrap

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/application/shared/auth"
	"github.com/HiroLiang/goat-server/internal/application/shared/security"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/database"
	dbUser "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/user"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis"
	infraAuth "github.com/HiroLiang/goat-server/internal/infrastructure/shared/auth"
	infraSecurity "github.com/HiroLiang/goat-server/internal/infrastructure/shared/security"
)

type Dependencies struct {
	Hasher       security.Hasher
	HMACer       security.HMACer
	TokenService auth.TokenService
	UserRepo     user.Repository
}

func BuildDeps() (*Dependencies, error) {

	// get config
	conf := config.App()

	// hmac configs
	hmacSecret := conf.Secrets.HmacSecret

	// AuthToken configs
	authExpiration := conf.AuthToken.Expiration

	return &Dependencies{
		Hasher:       infraSecurity.NewArgon2Hasher(),
		HMACer:       infraSecurity.NewSHA256HMACer(hmacSecret),
		TokenService: infraAuth.NewAuthTokenService(redis.RedisClient, time.Duration(authExpiration)*time.Second),
		UserRepo:     dbUser.NewUserRepository(database.Postgres),
	}, nil
}
