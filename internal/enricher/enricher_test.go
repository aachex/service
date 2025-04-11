package enricher

import (
	"fmt"
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
	var enriched model.EnrichedData

	err := EnrichUser(user, &enriched)
	if err != nil {
		t.Error(err)
	}
}

func TestEnrichAge(t *testing.T) {
	var enriched model.EnrichedData

	err := EnrichAge(user, &enriched)
	if err != nil {
		t.Error(err)
	}
}

func TestEnrichGender(t *testing.T) {
	var enriched model.EnrichedData

	err := EnrichGender(user, &enriched)
	if err != nil {
		t.Error(err)
	}

	if enriched.Gender == "" {
		t.Error("gender is empty")
	}
}

func TestEnrichNationality(t *testing.T) {
	var enriched model.EnrichedData

	err := EnrichNationality(user, &enriched)
	if err != nil {
		t.Error(err)
	}

	if enriched.Nationality == "" {
		t.Error("nationality is empty")
	}

	fmt.Println(enriched)
}
