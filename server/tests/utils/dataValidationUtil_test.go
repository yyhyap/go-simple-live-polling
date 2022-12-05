package tests_utils

import (
	"go-simple-live-polling/utils"
	"testing"
)

var (
	dataValidationUtil utils.IDataValidationUtil = utils.GetDataValidationUtil()
)

func TestValidIcNo(t *testing.T) {
	icNo := "961221049552"
	match := dataValidationUtil.IsIcNoValid(icNo)

	if !match {
		t.Error("should be true")
	}
}

func TestInvalidIcNo(t *testing.T) {
	icNo := "abcdefghijkl"
	match := dataValidationUtil.IsIcNoValid(icNo)

	if match {
		t.Error("should be false")
	}

	icNo = "961221049552 "
	match = dataValidationUtil.IsIcNoValid(icNo)

	if match {
		t.Error("should be false")
	}

	icNo = " 961221049552"
	match = dataValidationUtil.IsIcNoValid(icNo)

	if match {
		t.Error("should be false")
	}

	icNo = "96122104955a"
	match = dataValidationUtil.IsIcNoValid(icNo)

	if match {
		t.Error("should be false")
	}

	icNo = "9612210495529"
	match = dataValidationUtil.IsIcNoValid(icNo)

	if match {
		t.Error("should be false")
	}

	icNo = "96122104955"
	match = dataValidationUtil.IsIcNoValid(icNo)

	if match {
		t.Error("should be false")
	}
}
