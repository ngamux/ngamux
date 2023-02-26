package ngamux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-must/must"
)

func TestNewCtx(t *testing.T) {
	must := must.New(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	expected := &Ctx{
		Res(w),
		Req(r),
	}

	result := NewCtx(w, r)

	must.NotNil(result)
	must.Equal(expected, result)
}
