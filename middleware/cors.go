package middleware

import "strings"

type CorsOptions struct {
	AllowMethods      []string
	AllowHeaders      []string
	AllowOrigins      []string
	AllowExposeHeader []string
}

type Cors struct {
	Origin       string
	Method       string
	Header       string
	ExposeHeader string
}

func EnableCors(options CorsOptions) *Cors {

	c := new(Cors)

	if len(options.AllowHeaders) == 0 || options.AllowHeaders[0] == "*" {
		c.Header = "*"
	} else {
		c.Header = strings.Join(options.AllowHeaders, ",")
	}

	if len(options.AllowOrigins) == 0 || options.AllowOrigins[0] == "*" {
		c.Origin = "*"
	} else {
		c.Origin = strings.Join(options.AllowOrigins, ",")
	}

	if len(options.AllowExposeHeader) == 0 || options.AllowExposeHeader[0] == "*" {
		c.ExposeHeader = "*"
	} else {
		c.ExposeHeader = strings.Join(options.AllowExposeHeader, ",")
	}

	c.Method = strings.Join(options.AllowMethods, ",")

	return c
}
