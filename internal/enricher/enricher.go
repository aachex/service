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

func EnrichAge(user model.User, enriched *model.EnrichedUser) error {
	type resBody struct {
		Age int `json:"age"`
	}

	body, err := httpGet[resBody]("https://api.agify.io/?name=" + user.Name)
	if err != nil {
		return err
	}

	enriched.Age = body.Age
	return nil
}

func EnrichGender(user model.User, enriched *model.EnrichedUser) error {
	type resBody struct {
		Gender string `json:"gender"`
	}

	body, err := httpGet[resBody]("https://api.genderize.io/?name=" + user.Name)
	if err != nil {
		return err
	}

	enriched.Gender = body.Gender
	return nil
}

func httpGet[T any](url string) (body T, err error) {
	res, err := http.Get(url)
	if err != nil {
		return body, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return body, err
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		return body, err
	}
	return body, nil
}
