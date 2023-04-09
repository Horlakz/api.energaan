package helper

import (
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

func UuidPointerToString(uid *uuid.UUID) (str string) {
	if uid == nil {
		return ""
	}

	return uid.String()
}

func StringToUuidPointer(str string) (uid *uuid.UUID) {
	if str == "00000000-0000-0000-0000-000000000000" || str == "" {
		return nil
	}

	nuid, _ := uuid.Parse(str)
	return &nuid
}

func CreateSlug(str string) string {
	text := slug.Make(str)

	return text
}
