package mydb

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lastforeverzl/barkme-server/config"
	"github.com/lastforeverzl/barkme-server/message"
)

type Datastore interface {
	GetAllUsers(chan *AllUsers)
	GetTokens(chan *TokensChan)
	CreateUser(chan *UserChan)
	UpdateUser(chan *UserChan, string, User)
	AddFavUser(chan *UserChan, string, User)
	RemoveFavUser(chan *UserChan, string, User)
	UpdateUserAction(message.Envelope)
	GetTokensPip() <-chan *TokensChan
}

type DB struct {
	*gorm.DB
}

func NewDB(cfg *config.ServerConfig) (*DB, error) {
	dbInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port)
	db, err := gorm.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	if err = db.DB().Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) InitSchema() {
	db.AutoMigrate(&User{})
}
