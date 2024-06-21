package model

import (
	"encoding/hex"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooLong                      = errors.New("password too long")
	ErrPasswordContainsUnsuportedCharacters = errors.New("password contains unsuported characters")
	ErrInvalidUsername                      = errors.New("invalid username")
	AUTH_BASE_USER                          = Authority{AuthName: "User", Level: 0}
	AUTH_MODERATOR                          = Authority{AuthName: "Moderator", Level: 1}
	AUTH_ADMIN                              = Authority{AuthName: "Admin", Level: 255}
)

type User struct {
	ID         uint64
	Username   string   `gorm:"unique"`
	Password   Password `gorm:"embedded"`
	Profile    Profile  `gorm:"embedded"`
	Email      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	FullName   string
	FollowList FollowList `gorm:"foreignKey:Owner;references:Username"`
	Active     bool       `gorm:"default:true"`
	Authority
}

type FollowList struct {
	ID        uint64
	Owner     string `gorm:"unique"`
	Following []User `gorm:"many2many:follows;foreignKey:Owner;joinForeignKey:Owner;joinReferences:Username;references:Username"`
}
type Password struct {
	HashedPassword string
	UpdatedAt      time.Time
}

type Profile struct {
	Bio          string
	PfPUrl       string
	PfPDeleteUrl string
}

type Section struct {
	ID    uint64
	Name  string
	Owner string
	User  User   `gorm:"foreignKey:Owner;references:Username"`
	Posts []Post `gorm:"many2many:section_posts;"`
}

type Authority struct {
	AuthName string
	Level    uint8
}

func NewUser() User {
	return User{
		Authority: AUTH_BASE_USER,
		Active:    true,
		Profile:   Profile{PfPUrl: "/static/default-avatar.png"},
	}
}

func (p *Password) SetPasswordAsHash(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	encodedHashedPassword := hex.EncodeToString(hashedPassword)
	p.HashedPassword = encodedHashedPassword
	return nil
}

func (p *Password) ComparePassword(password string) bool {
	decodedHashedPassword, err := hex.DecodeString(p.HashedPassword)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword(decodedHashedPassword, []byte(password))
	return err == nil
}

func (p *Password) ValidateAndSetPassword(password string) error {
	passwordLength := len(password)
	regex := regexp.MustCompile(`^[a-zA-Z0-9!#$%&()*+,-.:;<=>?@[\]_{} ]+$`)
	ok := regex.MatchString(password)
	if !ok {
		return ErrPasswordContainsUnsuportedCharacters
	}
	if passwordLength > 72 || passwordLength < 12 {
		return ErrPasswordTooLong
	}
	return p.SetPasswordAsHash(password)
}

func isValidUsername(username string) bool {
	// Expresión regular para verificar el formato del nombre de usuario
	if len(username) < 5 || len(username) > 20 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return re.MatchString(username)
}

func (u *User) ValidateAndSetUsername(username string) error {
	if !isValidUsername(username) {
		return ErrInvalidUsername
	}
	u.Username = username
	return nil
}
