package utils

import (
	"github.com/avialog/backend/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"slices"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	err := validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		role := fl.Field().String()
		return slices.Contains(model.AvailableRoles, model.Role(role))
	})
	if err != nil {
		logrus.Panic(err)
	}

	err = validate.RegisterValidation("approach_type", func(fl validator.FieldLevel) bool {
		approachType := fl.Field().String()
		return slices.Contains(model.AvailableApproachTypes, model.ApproachType(approachType))
	})
	if err != nil {
		logrus.Panic(err)
	}

	err = validate.RegisterValidation("style", func(fl validator.FieldLevel) bool {
		style := fl.Field().String()
		return slices.Contains(model.AvailableStyles, model.Style(style))
	})
	if err != nil {
		logrus.Panic(err)
	}
}

func GetValidator() *validator.Validate {
	return validate
}
