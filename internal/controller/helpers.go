package controller

import (
	"encoding/json"
	"io"
	"net/http"
)

// readBody получает из тела запроса в формате json структуру T.
func readBody[T any](r *http.Request) (*T, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var obj T
	err = json.Unmarshal(bodyBytes, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

// writeReponse записывает структуру T в ответ в формате json.
func writeReponse[T any](obj T, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(&obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(b)
}
