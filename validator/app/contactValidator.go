package app

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
	appDto "github.com/horlakz/energaan-api/database/dto/app"
	"github.com/horlakz/energaan-api/validator"
)

type ContactValidator struct {
	validator.Validator[appDto.ContactDTO]
}

func (validator *ContactValidator) Validate(contactDto appDto.ContactDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&contactDto,
		validation.Field(&contactDto.FullName, validation.Required, validation.Length(1, 32)),
		validation.Field(&contactDto.Email, validation.Required, validation.Length(1, 32)),
		validation.Field(&contactDto.Phone, validation.Required),
		validation.Field(&contactDto.Country, validation.Required),
		validation.Field(&contactDto.Message, validation.Required),
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
