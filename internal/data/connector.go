package data

import (
	"context"
	"github.com/IgorViskov/33_loyalty/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Connector interface {
	GetConnection(context context.Context) *gorm.DB
	Close() error
}

type connector struct {
	db    *gorm.DB
	dbURI string
	mutex sync.Mutex
}

func NewConnector(conf *config.AppConfig) Connector {
	return &connector{
		dbURI: conf.DbURI,
	}
}

func (c *connector) GetConnection(context context.Context) *gorm.DB {
	if c.db == nil {
		c.mutex.Lock()
		if c.db == nil {
			var err error
			c.db, err = c.connect()
			if err != nil {
				c.mutex.Unlock()
				panic(err)
			}
		}
		c.mutex.Unlock()
	}
	return c.db.Session(&gorm.Session{
		Context: context,
	})
}

func (c *connector) Close() error {
	db, err := c.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (c *connector) connect() (*gorm.DB, error) {
	_, err := pgxpool.ParseConfig(c.dbURI)
	if err != nil {
		return nil, err
	}
	return gorm.Open(postgres.Open(c.dbURI))
}
