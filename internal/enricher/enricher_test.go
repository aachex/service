package enricher

import (
	"testing"

	"github.com/aachex/service/internal/model"
)

var user = model.User{
	Id:         921,
	Name:       "Dmitry",
	Surname:    "Dimov",
	Patronymic: "Dmitrievich",
}

func TestEnrichUser(t *testing.T) {
	err := EnrichUser(&user)
	if err != nil {
		t.Error(err)
	}
}

func TestEnrichAge(t *testing.T) {
	err := EnrichAge(&user)
	if err != nil {
		t.Error(err)
	}
}

func TestEnrichGender(t *testing.T) {
	err := EnrichGender(&user)
	if err != nil {
		t.Error(err)
	}

	if user.Gender == "" {
		t.Error("gender is empty")
	}
}

func TestEnrichNationality(t *testing.T) {
	err := EnrichNationality(&user)
	if err != nil {
		t.Error(err)
	}

	if user.Nationality == "" {
		t.Error("nationality is empty")
	}
}
