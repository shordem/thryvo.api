package database

import (
	"gorm.io/gorm"

	"github.com/shordem/api.thryvo/lib/constants"
)

type DatabaseInterface interface {
	Connection() *gorm.DB
	Cache() RedisClientInterface
}

type connection struct {
	pg    PostgresClientInterface
	cache RedisClientInterface
}

func StartDatabaseClient(env constants.Env) DatabaseInterface {
	return &connection{
		pg:    NewPostgresClient(env),
		cache: NewRedisClient(env),
	}
}

func (conn connection) Connection() *gorm.DB {
	return conn.pg.Connection()
}

func (conn connection) Cache() RedisClientInterface {
	return conn.cache
}
