package ngamux

import "net/http"

type Ctx struct {
	Res *Response
	Req *Request
}

func NewCtx(w http.ResponseWriter, r *http.Request) *Ctx {
	return &Ctx{
		Res(w),
		Req(r),
	}
}
