package ngamux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Add("content-type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Fprint(rw, string(jsonData))
	return nil
}

func JSONWithStatus(rw http.ResponseWriter, status int, data interface{}) error {
	rw.WriteHeader(status)
	err := JSON(rw, data)
	if err != nil {
		return err
	}

	return nil
}

func GetParam(r *http.Request, key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

func GetBody(r *http.Request, store interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&store); err != nil {
		return err
	}

	return nil
}
