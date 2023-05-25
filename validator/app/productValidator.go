package app

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
	appDto "github.com/horlakz/energaan-api/database/dto/app"
	"github.com/horlakz/energaan-api/validator"
)

type ProductValidator struct {
	validator.Validator[appDto.ProductDTO]
}

func (validator *ProductValidator) Validate(productDto appDto.ProductDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&productDto,
		validation.Field(&productDto.Slug,
			validation.Required,
			validation.Length(1, 32)),
		validation.Field(&productDto.Title, validation.Required, validation.Length(1, 32)),
		validation.Field(&productDto.CategoryID, validation.Required),
		validation.Field(&productDto.Images, validation.Required, validation.Length(1, 64)),
		validation.Field(&productDto.Description, validation.Required),
		validation.Field(&productDto.Features, validation.Required),
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
