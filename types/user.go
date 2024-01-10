package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	MinFirstNameLen = 2
	MinLastNameLen  = 2
	minPasswordLen  = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type UpdateUserParams struct {
	ID        primitive.ObjectID
	FirstName string
	LastName  string
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < MinFirstNameLen {
		errors["firstNameError"] = fmt.Sprintf("first name must have atleast %d characters :)", MinFirstNameLen)
	}
	if len(params.LastName) < MinLastNameLen {
		errors["lastNameError"] = fmt.Sprintf("last name must have atleast %d characters :)", MinLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["passwordError"] = fmt.Sprintf("password must have atleast %d characters :)", minPasswordLen)
	}
	if !params.isEmailValid() {
		errors["emailError"] = "email address is invalid :("
	}
	return errors

}
func (params CreateUserParams) isEmailValid() bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.Match([]byte(params.Email))
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
