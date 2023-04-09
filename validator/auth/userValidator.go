package auth

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"

	authDtos "github.com/horlakz/energaan-api/database/dto/auth"
	"github.com/horlakz/energaan-api/validator"
)

type UserValidator struct {
	validator.Validator[authDtos.UserDTO]
}

func (validator *UserValidator) Validate(userDto authDtos.UserDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&userDto,
		validation.Field(&userDto.FullName, validation.Required, validation.Length(6, 128)),
		validation.Field(&userDto.Email, validation.Required, validation.Length(2, 128), validation.By(validator.ValidateDBUnique(userDto, "admins", "email", map[string]interface{}{}))),
		validation.Field(&userDto.Password, validation.Required, validation.Length(6, 128)),
	)

	if err != nil {
		if e, ok := err.(validation.InternalError); ok {
			log.Println(e.InternalError())
			return nil, nil
		}

		var dat map[string]interface{}
		m, _ := json.Marshal(err)

		json.Unmarshal(m, &dat)
		return dat, err
	}

	return nil, nil
}
