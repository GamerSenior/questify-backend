package models

import (
	"github.com/dgrijalva/jwt-go"
	u "github.com/GamerSenior/questify-back/utils"
	"strings"
	"github.com/jinzhu/gorm"
	"os"
	"golang.org/x/crypto/bcrypt"
)

// Token struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account model - used for login
type Account struct {
	gorm.Model
	Email string `json:"email" gorm:"unique not null"`
	Password string `json:"password"`
	Token string `json:"token" sql:"-"`
}

// Validate validates the Account
func (account *Account) Validate() (map[string] interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password must contain at least 6 characters"), false
	}

	temp := &Account{}

	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Error connecting database, please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already registered"), false
	}

	return u.Message(false, "Requirement passed"), true
}

// Create creates new Account entity
func (account *Account) Create() (map[string] interface{}, bool) {
	if resp, ok := account.Validate(); !ok {
		return resp, false
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost) 
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Faile to create account, please retry"), false
	}

	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = ""

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response, true
}

// Login validates account passwords and sets Token
func Login(email, password string) (map[string] interface{}) {
	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error. Please retr")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)) 
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { 
		return u.Message(false, "Invalid login credentials. Please try again")
	} 

	account.Password = ""

	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	resp := u.Message(true, "Logged in")
	resp["account"] = account
	return resp
}

// GetUser returns the Account for given ID
func GetUser(u uint) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}
	acc.Password = ""
	return acc
}