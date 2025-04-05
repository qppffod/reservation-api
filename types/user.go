package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	minFirstNameLen    = 2
	minLastNameLen     = 2
	minPasswordNameLen = 7
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}

	return m
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}

	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordNameLen {
		errors["password"] = fmt.Sprintf("password should be at least %d characters", minPasswordNameLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = "invalid email"
	}

	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(user *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
