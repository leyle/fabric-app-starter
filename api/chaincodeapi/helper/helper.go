package helper

import (
	"github.com/google/uuid"
	"github.com/leyle/go-api-starter/logmiddleware"
)

func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return true
}

func GeneratePrivateDataId() string {
	return uuid.New().String()
}

func GenerateTransientKey() string {
	return logmiddleware.GenerateReqId()
}
