package enricher

import (
	"testing"

	"github.com/aachex/service/internal/model"
)

var user = model.User{
	Id:         921,
	Name:       "Dima",
	Surname:    "Dimov",
	Patronymic: "Dmitrievich",
}

func TestEnrichAge(t *testing.T) {
	var enriched model.EnrichedUser
	enriched.User = user

	err := EnrichAge(user, &enriched)
	if err != nil {
		t.Error(err)
	}
}

func TestEnrichGender(t *testing.T) {
	var enriched model.EnrichedUser
	enriched.User = user

	err := EnrichGender(user, &enriched)
	if err != nil {
		t.Error(err)
	}

	if enriched.Gender == "" {
		t.Error("gender is empty")
	}
}
