package enricher

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aachex/service/internal/model"
)

func EnrichUser(user model.User) (model.EnrichedUser, error) {
	return model.EnrichedUser{}, nil
}

func EnrichAge(user model.User) (result model.EnrichedUser, err error) {
	result.User = user
	url := "https://api.agify.io/?name=" + user.Name

	res, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return result, err
	}

	var body struct {
		Count int `json:"count"`
		Age   int `json:"age"`
	}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return result, err
	}
	result.Age = body.Age
	return result, nil
}
