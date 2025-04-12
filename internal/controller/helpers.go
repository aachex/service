package controller

import (
	"encoding/json"
	"io"
	"net/http"
)

// readBody получает из тела запроса в формате json структуру T.
func readBody[T any](r *http.Request) (obj T, err error) {
	if r.Body == nil {
		return obj, nil
	}

	err = json.NewDecoder(r.Body).Decode(&obj)
	if err == io.EOF {
		return obj, nil
	}
	if err != nil {
		return obj, err
	}

	return obj, nil
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
