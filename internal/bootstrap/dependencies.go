package bootstrap

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/application/shared/auth"
	"github.com/HiroLiang/goat-server/internal/application/shared/security"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/domain/agent"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/domain/userrole"
	"github.com/HiroLiang/goat-server/internal/infrastructure/auth/session"
	infraAuth "github.com/HiroLiang/goat-server/internal/infrastructure/auth/token"
	"github.com/HiroLiang/goat-server/internal/infrastructure/cache"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/database"
	dbAgent "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/agent"
	dbUser "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/user"
	dbUserrole "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/userrole"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis"
	userrole2 "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis/userrole"
	infraSecurity "github.com/HiroLiang/goat-server/internal/infrastructure/shared/security"
)

type Dependencies struct {
	Hasher       security.Hasher
	HMACer       security.HMACer
	TokenService auth.TokenService
	UserRepo     user.Repository
	UserRoleRepo userrole.Repository
	AgentRepo    agent.Repository
}

func BuildDeps() (*Dependencies, error) {

	// get config
	conf := config.App()

	// redis cache
	var redisCache cache.Cache = redis.NewRedisCache(redis.RedisClient)

	// hmac configs
	hmacSecret := conf.Secrets.HmacSecret

	// AuthToken configs
	authExpiration := conf.AuthToken.Expiration

	return &Dependencies{
		Hasher: infraSecurity.NewArgon2Hasher(),
		HMACer: infraSecurity.NewSHA256HMACer(hmacSecret),
		TokenService: infraAuth.NewAuthTokenService(
			session.NewRedisSessionStore(redisCache, redis.RedisClient),
			time.Duration(authExpiration)*time.Second),
		UserRepo:  dbUser.NewUserRepository(database.Postgres),
		AgentRepo: dbAgent.NewAgentRepository(database.Postgres),
		UserRoleRepo: userrole2.NewUserRoleCachedRepo(redisCache,
			dbUserrole.NewUserRoleRepository(database.Postgres)),
	}, nil
}
