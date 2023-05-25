package app

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
	appDto "github.com/horlakz/energaan-api/database/dto/app"
	"github.com/horlakz/energaan-api/validator"
)

type QuoteValidator struct {
	validator.Validator[appDto.QuoteDTO]
}

func (validator *QuoteValidator) Validate(quoteDto appDto.QuoteDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&quoteDto,
		validation.Field(&quoteDto.FullName, validation.Required, validation.Length(1, 32)),
		validation.Field(&quoteDto.Email, validation.Required, validation.Length(1, 32)),
		validation.Field(&quoteDto.ServiceId, validation.Required),
		validation.Field(&quoteDto.ServiceType, validation.Required),
		validation.Field(&quoteDto.Phone, validation.Required),
		validation.Field(&quoteDto.Country, validation.Required),
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
