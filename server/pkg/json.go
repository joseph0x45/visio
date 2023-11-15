package pkg

import (
	"encoding/json"
	"net/http"
)

func RespondToBadRequest(w http.ResponseWriter, error_code string) error {
	response, err := json.Marshal(
    map[string]string{
      "error": error_code,
    },
  )
	if err != nil {
		return err
	}
  w.WriteHeader(http.StatusBadRequest)
  w.Write(response)
  return nil
}
