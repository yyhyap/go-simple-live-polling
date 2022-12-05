package utils

import (
	"go-simple-live-polling/logger"
	"regexp"
	"sync"
)

var (
	dataValidationUtilOnce sync.Once
	dataValidationUtil     *dataValidationUtilStruct
)

type IDataValidationUtil interface {
	IsIcNoValid(s string) bool
}

type dataValidationUtilStruct struct{}

func GetDataValidationUtil() *dataValidationUtilStruct {
	if dataValidationUtil == nil {
		dataValidationUtilOnce.Do(func() {
			dataValidationUtil = &dataValidationUtilStruct{}
		})
	}
	return dataValidationUtil
}

func (d *dataValidationUtilStruct) IsIcNoValid(s string) bool {
	match, err := regexp.MatchString(`^(\d{12})$`, s)
	if err != nil {
		logger.Logger.Error(err.Error())
		return false
	}
	return match
}
