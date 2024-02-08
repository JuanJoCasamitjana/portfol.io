package model

import (
	"encoding/hex"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooLong = errors.New("password too long")
var ErrPasswordContainsUnsuportedCharacters = errors.New("password contains unsuported characters")
var ErrInvalidUsername = errors.New("invalid username")

// Migration should be Password -> User -> Auth
// For now no email will be required
// Pointers allow for null values
type User struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `gorm:"unique"`
	Firstname *string
	Lastname  *string
	AuthID    uint64   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Password  Password `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Auth struct {
	ID    uint64 `gorm:"primaryKey"`
	Level uint8
	Users []User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Password struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	Hash      string //saved as string to ensure compatibility
}

// Password utilities

func (p *Password) SetPasswordAsHash(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	encodedHashedPassword := hex.EncodeToString(hashedPassword)
	p.Hash = encodedHashedPassword
	return nil
}

func (p *Password) ComparePassword(password string) bool {
	decodedHashedPassword, err := hex.DecodeString(p.Hash)
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
	if passwordLength > 72 {
		return ErrPasswordTooLong
	}
	return p.SetPasswordAsHash(password)
}

func isValidUsername(username string) bool {
	// Expresión regular para verificar el formato del nombre de usuario
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	return re.MatchString(username)
}

func (u *User) ValidateAndSetUsername(username string) error {
	if !isValidUsername(username) {
		return ErrInvalidUsername
	}
	u.Username = username
	return nil
}
