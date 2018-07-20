package mydb

import (
	"fmt"

	"github.com/lastforeverzl/barkme-server/message"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	DeviceName string  `json:"deviceName"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
	Barks      int32   `json:"barks"`
	Token      string  `json:"token"`
	Favorites  []*User `gorm:"many2many:friendships;association_jointable_foreignkey:friend_id"`
}

type AllUsers struct {
	Users []User
	Err   error
}

type UserChan struct {
	User User
	Err  error
}

type TokensChan struct {
	Tokens []string
	Err    error
}

func (db *DB) GetAllUsers(c chan *AllUsers) {
	users := make([]User, 0)
	if err := db.Find(&users).Error; err != nil {
		c <- &AllUsers{Err: err}
	}
	c <- &AllUsers{Users: users}
	close(c)
}

func (db *DB) CreateUser(c chan *UserChan) {
	user := User{}
	if err := db.Create(&user).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	c <- &UserChan{User: user}
	close(c)
}

func (db *DB) UpdateUser(c chan *UserChan, id string, userUpdate User) {
	user := User{}
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	if err := db.Model(&user).Updates(userUpdate).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	c <- &UserChan{User: user}
	close(c)
}

func (db *DB) AddFavUser(c chan *UserChan, id string, favoriteUser User) {
	user := User{}
	favUser := User{}
	db.Preload("Favorites").First(&user, "id = ?", id)
	if err := db.Where("id = ?", favoriteUser.ID).First(&favUser).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	if err := db.Model(&user).Association("Favorites").Append(favUser).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	c <- &UserChan{User: user}
	close(c)
}

func (db *DB) RemoveFavUser(c chan *UserChan, id string, rmFavUser User) {
	user := User{}
	rmUser := User{}
	db.Preload("Favorites").First(&user, "id = ?", id)
	if err := db.Where("id = ?", rmFavUser.ID).First(&rmUser).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	if err := db.Model(&user).Association("Favorites").Delete(rmUser).Error; err != nil {
		c <- &UserChan{Err: err}
	}
	c <- &UserChan{User: user}
	close(c)
}

func (user *User) AfterCreate(tx *gorm.DB) (err error) {
	tx.Model(user).Update("DeviceName", fmt.Sprintf("#%d", user.ID))
	return
}

func (db *DB) UpdateUserAction(msg message.Envelope) {
	user := User{}
	db.Where("id = ?", msg.ID).First(&user)
	switch msg.Msg {
	case "barks":
		user.Barks++
	}
	db.Save(&user)
	fmt.Println(user)
}

func (db *DB) GetTokens(c chan *TokensChan) {
	users := make([]User, 0)
	if err := db.Select("token").Find(&users).Error; err != nil {
		c <- &TokensChan{Err: err}
	}
	tokens := make([]string, 0)
	for _, user := range users {
		if user.Token != "" {
			tokens = append(tokens, user.Token)
		}
	}
	c <- &TokensChan{Tokens: tokens}
	close(c)
}

func (db *DB) GetTokensPip() <-chan *TokensChan {
	out := make(chan *TokensChan)
	users := make([]User, 0)
	go func() {
		if err := db.Select("token").Find(&users).Error; err != nil {
			out <- &TokensChan{Err: err}
		}
		tokens := make([]string, 0)
		for _, user := range users {
			if user.Token != "" {
				tokens = append(tokens, user.Token)
			}
		}
		out <- &TokensChan{Tokens: tokens}
		close(out)
	}()
	return out
}
