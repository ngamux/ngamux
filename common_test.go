package ngamux

import (
	"net/http"
	"reflect"
	"testing"
)

func TestWithMiddlewares(t *testing.T) {
	result := WithMiddlewares()(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	if result == nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", reflect.TypeOf(result), nil)
	}

	result = WithMiddlewares(nil)(func(rw http.ResponseWriter, r *http.Request) error {
		return nil
	})
	if result == nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", reflect.TypeOf(result), nil)
	}

	result = WithMiddlewares(nil)(nil)
	if result != nil {
		t.Errorf("TestWithMiddlewares need %v, but got %v", nil, reflect.TypeOf(result))
	}
}
