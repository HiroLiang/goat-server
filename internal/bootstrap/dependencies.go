package bootstrap

import (
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
	redisUserrole "github.com/HiroLiang/goat-server/internal/infrastructure/persistence/redis/userrole"
	infraSecurity "github.com/HiroLiang/goat-server/internal/infrastructure/shared/security"
	"github.com/redis/go-redis/v9"
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

	// Postgres datasource
	postgres := dataSources.GetDB(database.Postgres)

	// redis cache
	redisCache := redisInfra.NewRedisCache(redis)

	// Session store
	sessionStore := session.NewRedisSessionStore(redisCache, redis)

	return &Dependencies{
		AgentRepo:    dbAgent.NewAgentRepository(postgres),
		TokenService: infraAuth.NewAuthTokenService(sessionStore, conf.AuthToken.Expiration),
		Hasher:       infraSecurity.NewArgon2Hasher(),
		HMACer:       infraSecurity.NewSHA256HMACer(conf.Secrets.HmacSecret),
		RateLimiter:  buildRateLimiter(redis, conf),
		UserRepo:     dbUser.NewUserRepository(postgres),
		UserRoleRepo: redisUserrole.NewUserRoleCachedRepo(redisCache, dbUserrole.NewUserRoleRepository(postgres)),
	}, nil
}

func BuildMockDeps(opts ...DepsOption) *Dependencies {
	conf := config.App()

	//Default dependencies
	deps := &Dependencies{
		Hasher: infraSecurity.NewArgon2Hasher(),
		HMACer: infraSecurity.NewSHA256HMACer(conf.Secrets.HmacSecret),
	}

	// Optionals
	for _, opt := range opts {
		opt(deps)
	}

	return deps
}

type DepsOption func(*Dependencies)

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
	return infraSecurity.NewRedisRateLimiter(
		redisInfraSecurity.NewRedisRateLimitRepository(redis),
		globalPolicy,
		ipPolicy)
}
