package users

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"donategold.me/db"
	"github.com/tarantool/go-tarantool"
	"gopkg.in/kataras/iris.v6"
)

const tableNameAT string = "access_tokens"

func init() {
	createTablesTarantool()
	db.DB.AutoMigrate(&User{}, &Provider{})
}

func createTablesTarantool() {
	qCreateTable := fmt.Sprintf("box.schema.space.create('%s', {if_not_exists=true})", tableNameAT)
	var err error
	_, err = db.TDB.Eval(qCreateTable, []interface{}{})
	if err != nil {
		log.Fatalln(err)
	}
	qCreateIndex := fmt.Sprintf(`box.space.%s:create_index('primary', {type = 'hash', if_not_exists=true, parts = {1, 'string'}})`, tableNameAT)
	_, err = db.TDB.Eval(qCreateIndex, []interface{}{})
	if err != nil {
		log.Fatalln(err)
	}
}

const mySigningKey = "AllYourBase"

type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Username  string     `gorm:"not null;unique" json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func (u *User) Create() error {
	return db.DB.Create(&u).Error
}

func (u *User) CreateUniq() error {
	err := db.DB.Create(&u).Error
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
		var user User
		db.DB.Last(&user)
		u.Username += fmt.Sprintf("%d", user.ID)
		return u.Create()
	} else if err != nil {
		return err
	}
	return nil
}

func (u *User) Valid() error {
	if db.DB.Where("username = ?", u.Username).First(&u); u.ID != 0 {
		return errors.New("Пользователь с таким именем уже существует")
	}
	return nil
}

func getAllUsers() []User {
	var users []User
	db.DB.Find(&users)
	return users
}

func GetUserById(id uint) (User, bool) {
	var u User
	var exist bool
	db.DB.Where("id = ?", id).First(&u)
	if u.ID != 0 {
		exist = true
	}
	return u, exist
}

func GetUserByUsername(username string) (User, bool) {
	var u User
	var exist bool
	db.DB.Where("username = ?", username).First(&u)
	if u.ID != 0 {
		exist = true
	}
	return u, exist
}

type AccessToken struct {
	Token string `json:"token"`
	User
}

func (at AccessToken) Create() error {
	_, err := db.TDB.Replace(tableNameAT, at)
	if err != nil {
		return err
	}
	return nil
}

func (at AccessToken) Delete() error {
	_, err := db.TDB.Delete(tableNameAT, "primary", []interface{}{at.Token})
	return err
}

func (at AccessToken) IsAuthenticated() bool {
	if at.Username == "" {
		return false
	}
	return true
}

func NewTokenForUser(user User) AccessToken {
	t := time.Now()
	token := fmt.Sprintf("%x", md5.Sum([]byte(user.Username+user.Email+t.String())))
	at := AccessToken{Token: token, User: user}
	err := at.Create()
	if err != nil {
		log.Println(err)
		return AccessToken{}
	}
	return at
}

func GetToken(token string) (AccessToken, bool) {
	var tokens = make([]AccessToken, 0)
	err := db.TDB.SelectTyped(tableNameAT, "primary", 0, 1, tarantool.IterEq, []interface{}{token}, &tokens)
	if err != nil || len(tokens) < 1 {
		return AccessToken{}, false
	}
	return tokens[0], true
}

type Provider struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	Uid          string     `sql:"index" json:"uid"`
	UserID       uint       `sql:"index" json:"user_id"`
	TypeProvider string     `sql:"index" json:"type_provider"`
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	Expires      *time.Time `json:"expires"`
}

func (p *Provider) Create() {
	db.DB.Where("uid = ? and type_provider = ?", p.Uid, p.TypeProvider).Delete(Provider{})
	db.DB.Create(p)
}

func (p *Provider) Save() {
	db.DB.Save(p)
}

func (p *Provider) Delete() {
	db.DB.Delete(p)
}

func (p *Provider) RefreshTokenNeed() bool {
	if p.Expires == nil {
		return false
	}
	return p.Expires.Before(time.Now())
}

func GetProvidersByUser(userID uint) []Provider {
	var providers []Provider
	db.DB.Where("user_id = ?", userID).Find(&providers)
	return providers
}

func GetProviderByTypeAndUserID(typeOfP string, userID uint) (Provider, error) {
	if typeOfP == "youtube" {
		typeOfP = "gplus"
	}
	var provider Provider
	db.DB.Where("type_provider = ? and user_id = ?", typeOfP, userID).First(&provider)
	if provider.ID == 0 {
		return provider, errors.New("Provider not found")
	}
	return provider, nil
}

func GetProviderByTypeAndUserIDAndUID(typeOfP string, userID uint, uid string) (Provider, error) {
	if typeOfP == "youtube" {
		typeOfP = "gplus"
	}
	var provider Provider
	db.DB.Where("type_provider = ? and user_id = ? and uid = ?", typeOfP, userID, uid).First(&provider)
	if provider.ID == 0 {
		return provider, errors.New("Provider not found")
	}
	return provider, nil
}

func deleteProviderByIdAndUserId(userID uint, id uint) error {
	err := db.DB.Where("user_id = ? and id = ?", userID, id).Delete(Provider{}).Error
	return err
}

// @reqired@ - required login
func GetMiddleware(required bool) func(c *iris.Context) {
	return func(c *iris.Context) {
		tokenStr := c.RequestHeader("X-TOKEN-KEY")
		if tokenStr == "" {
			tokenStr = c.URLParam("token")
		}
		if tokenStr == "" && c.Method() != "GET" && required {
			c.EmitError(http.StatusUnauthorized)
			return
		}
		if tokenStr == "" {
			tokenStr = c.GetCookie("token")
		}
		token, exist := GetToken(tokenStr)
		if required && !exist {
			c.EmitError(http.StatusUnauthorized)
			return
		}
		c.Set("token", token)
		c.Next()
	}
}

func CentrifugoMidll(c *iris.Context) {
	tokenStr := c.GetCookie("token")
	token, exist := GetToken(tokenStr)
	if !exist {
		c.EmitError(http.StatusUnauthorized)
		return
	}
	c.Set("token", token)
	c.Next()
}
