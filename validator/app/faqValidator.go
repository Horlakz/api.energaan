package app

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
	appDto "github.com/horlakz/energaan-api/database/dto/app"
	"github.com/horlakz/energaan-api/validator"
)

type FaqValidator struct {
	validator.Validator[appDto.FaqDTO]
}

func (validator *FaqValidator) Validate(faqDto appDto.FaqDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&faqDto,
		validation.Field(&faqDto.Title, validation.Required, validation.Length(1, 32)),
		validation.Field(&faqDto.Description, validation.Required),
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
