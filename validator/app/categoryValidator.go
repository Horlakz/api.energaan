package app

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
	appDto "github.com/horlakz/energaan-api/database/dto/app"
	"github.com/horlakz/energaan-api/validator"
)

type CategoryValidator struct {
	validator.Validator[appDto.CategoryDTO]
}

func (validator *CategoryValidator) Validate(categoryDto appDto.CategoryDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&categoryDto,
		validation.Field(&categoryDto.Slug, validation.Required, validation.Length(1, 32)),
		validation.Field(&categoryDto.Name, validation.Required, validation.Length(1, 32)),
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
