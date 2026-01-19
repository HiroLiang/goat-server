package bootstrap

import (
	"time"

	"github.com/HiroLiang/goat-server/internal/application/shared/auth"
	"github.com/HiroLiang/goat-server/internal/application/shared/security"
	"github.com/HiroLiang/goat-server/internal/config"
	"github.com/HiroLiang/goat-server/internal/domain/agent"
	domainSecurity "github.com/HiroLiang/goat-server/internal/domain/security"
	"github.com/HiroLiang/goat-server/internal/domain/user"
	"github.com/HiroLiang/goat-server/internal/domain/userrole"
	"github.com/HiroLiang/goat-server/internal/infrastructure/auth/session"
	infraAuth "github.com/HiroLiang/goat-server/internal/infrastructure/auth/token"
	"github.com/HiroLiang/goat-server/internal/infrastructure/persistence/database"
	dbAgent "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/agent"
	dbUser "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/user"
	dbUserrole "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/postgres/userrole"
	redisInfra "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis"
	redisInfraSecurity "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis/security"
	userrole2 "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis/userrole"
	infraSecurity "github.com/HiroLiang/goat-server/internal/infrastructure/shared/security"
	"github.com/HiroLiang/goat-server/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Dependencies struct {
	AgentRepo    agent.Repository
	TokenService auth.TokenService
	Hasher       security.Hasher
	HMACer       security.HMACer
	RateLimiter  security.RateLimiter
	UserRepo     user.Repository
	UserRoleRepo userrole.Repository
}

func BuildDeps(redis *redis.Client, dataSources *database.DataSources) (*Dependencies, error) {

	// get config
	conf := config.App()

	// redis cache
	redisCache := redisInfra.NewRedisCache(redis)

	// hmac configs
	hmacSecret := conf.Secrets.HmacSecret

	// AuthToken configs
	authExpiration := conf.AuthToken.Expiration

	return &Dependencies{
		AgentRepo: dbAgent.NewAgentRepository(dataSources.GetDB(database.Postgres)),
		TokenService: infraAuth.NewAuthTokenService(
			session.NewRedisSessionStore(redisCache, redis),
			time.Duration(authExpiration)*time.Second),
		Hasher:      infraSecurity.NewArgon2Hasher(),
		HMACer:      infraSecurity.NewSHA256HMACer(hmacSecret),
		RateLimiter: buildRateLimiter(redis, conf),
		UserRepo:    dbUser.NewUserRepository(dataSources.GetDB(database.Postgres)),
		UserRoleRepo: userrole2.NewUserRoleCachedRepo(redisCache,
			dbUserrole.NewUserRoleRepository(dataSources.GetDB(database.Postgres))),
	}, nil
}

// buildRateLimiter build rate limiter
func buildRateLimiter(redis *redis.Client, conf *config.AppConfig) security.RateLimiter {
	rateLimitConf := conf.RateLimitConfig
	globalPolicy := domainSecurity.RateLimitPolicy{
		Limit:  int64(rateLimitConf.GlobalLimit),
		Window: rateLimitConf.GlobalUnit,
	}
	ipPolicy := domainSecurity.RateLimitPolicy{
		Limit:  int64(rateLimitConf.IPLimit),
		Window: rateLimitConf.IPUnit,
	}
	logger.Log.Info("policy updated", zap.Any("global", globalPolicy))
	logger.Log.Info("policy updated", zap.Any("ip", ipPolicy))
	return infraSecurity.NewRedisRateLimiter(
		redisInfraSecurity.NewRedisRateLimitRepository(redis),
		globalPolicy,
		ipPolicy)
}
