package Modals

import "gopkg.in/go-playground/validator.v9"


type UserCrud struct {
    Name      string   `json:"name"                validate:"required"          firestore:"name"`
    Email     string   `json:"email"               validate:"required,email"    firestore:"email"`
	Password  string   `json:"password"            validate:"required,gte=6"    firestore:"password"`
	UserID    string   `json:"user_id,omitempty"    validate:"omitempty"    firestore:"user_id"`
}

func (u *UserCrud) Validator() error {
	validate := validator.New()
	err := validate.Struct(u)
	return err
}



type LoginCrud struct {
    Email  string      `json:"email"         validate:"required,email"        firestore:"email"`
    Password  string   `json:"password"      validate:"required,gte=6"        firestore:"password"`
}

func (u *LoginCrud) Validator() error {
	validate := validator.New()
	err := validate.Struct(u)
	return err
}



type FormCrud struct {
	Email 	 string	`json:"email"               validate:"required"          firestore:"email"`
	Phone    string `json:"phone"               validate:"required"          firestore:"phone"`
	Deadline string `json:"deadline"            validate:"required"          firestore:"deadline"`
	Detail	 string	`json:"detail"              validate:"required"          firestore:"detail"`
	Lang 	 string `json:"lang"                validate:"required"          firestore:"lang"`
	Type     string `json:"type"                validate:"required"          firestore:"type"`
	FileUrl  string `json:"fileUrl,omitempty"   validate:"omitempty"         firestore:"fileUrl,omitempty"`
}


func (u *FormCrud) Validator() error {
	validate := validator.New()
	err := validate.Struct(u)
	return err
}
