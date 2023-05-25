package app

import (
	"encoding/json"
	"log"

	validation "github.com/go-ozzo/ozzo-validation"
	appDto "github.com/horlakz/energaan-api/database/dto/app"
	"github.com/horlakz/energaan-api/validator"
)

type GalleryValidator struct {
	validator.Validator[appDto.GalleryDTO]
}

func (validator *GalleryValidator) Validate(galleryDto appDto.GalleryDTO) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&galleryDto,
		validation.Field(&galleryDto.Image, validation.Required, validation.Length(1, 64)),
		validation.Field(&galleryDto.Title, validation.Required, validation.Length(1, 32)),
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
