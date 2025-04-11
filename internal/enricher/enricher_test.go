package enricher

import (
	"fmt"
	"testing"

	"github.com/aachex/service/internal/model"
)

func TestEnrichAge(t *testing.T) {
	user := model.User{
		Id:         921,
		Name:       "Artem",
		Surname:    "Fatahov",
		Patronymic: "Dmitrievich",
	}

	enriched, err := EnrichAge(user)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(enriched)
}
