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
