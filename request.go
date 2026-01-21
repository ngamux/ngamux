package ngamux

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/ngamux/ngamux/json"
)

var (
	unsupportedFieldType = errors.New("unsupported field type")
)

// Request define single request manager
type Request struct {
	*http.Request
}

// Req needs *http.Request and returns *Request object
func Req(r *http.Request) *Request {
	return &Request{r}
}

// Params returns parameter from url using a key
func (r Request) Params(key string) string {
	params := r.Context().Value(KeyContextParams).([][]string)
	for _, param := range params {
		if param[0] == key {
			return param[1]
		}
	}

	return ""
}

// Query returns data from query params using a key
func (r Request) Query(key string, fallback ...string) string {
	queries := r.URL.Query()
	query := queries.Get(key)
	if query == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}
	return query
}

// Parse query into struct
func (r Request) QueriesParser(data any) error {
	queryValues := r.URL.Query()

	val := reflect.ValueOf(data).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("query")
		if tag != "" {
			if queryValue, found := queryValues[tag]; found && len(queryValue) > 0 {
				fieldValue := val.Field(i)
				if fieldValue.CanSet() {
					if err := setFieldValue(fieldValue, queryValue[0]); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func setFieldValue(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() {
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		fieldValue.SetInt(int64(intValue))
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		fieldValue.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, fieldValue.Type().Bits())
		if err != nil {
			return err
		}
		fieldValue.SetFloat(floatValue)
	case reflect.Struct:
		nestedStruct := reflect.New(fieldValue.Type()).Elem()
		if err := parseNestedStruct(nestedStruct, value); err != nil {
			return err
		}
		fieldValue.Set(nestedStruct)
	default:
		return unsupportedFieldType
	}
	return nil
}

func parseNestedStruct(nestedStruct reflect.Value, value string) error {
	queryValues, err := url.ParseQuery(value)
	if err != nil {
		return err
	}

	typ := nestedStruct.Type()
	for i := 0; i < nestedStruct.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("query")
		if tag != "" {
			if queryValue, found := queryValues[tag]; found && len(queryValue) > 0 {
				fieldValue := nestedStruct.Field(i)
				if fieldValue.CanSet() {
					if err := setFieldValue(fieldValue, queryValue[0]); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// FormValue returns data from form using a key
func (r Request) FormValue(key string, fallback ...string) string {
	value := r.PostFormValue(key)
	if value == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return ""
	}

	return value
}

// FormFile returns file from form using a key
func (r Request) FormFile(key string, maxFileSize ...int64) (*multipart.FileHeader, error) {
	var maxFileSizeParsed int64 = 10 << 20
	if len(maxFileSize) > 0 {
		maxFileSizeParsed = maxFileSize[0]
	}

	err := r.ParseMultipartForm(maxFileSizeParsed)
	if err != nil {
		return nil, err
	}

	file, header, err := r.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return header, nil
}

// JSON get json data from request body and store to variable reference
func (r Request) JSON(store any) error {
	rBody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rBody, &store)
	if err != nil {
		return err
	}

	return nil
}

// Locals needs key and optional value
// It returns any if only key and no value given
// It insert value to context if key and value is given
func (r *Request) Locals(key any, value ...any) any {
	if len(value) <= 0 {
		return r.Context().Value(key)
	}
	r.Request = r.WithContext(context.WithValue(r.Context(), key, value[0]))
	return nil
}

// IsLocalhost returns true if hostname is localhost or 127.0.0.1
func (r *Request) IsLocalhost() bool {
	return strings.Contains(r.Host, "localhost") || strings.Contains(r.Host, "127.0.0.1")
}

// Get Client IP
func (r Request) GetIPAdress() string {

	var ipAddress string

	// X-Real-Ip - fetches first true IP (if the requests sits behind multiple NAT sources/load balancer
	ipAddress = r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		// X-Forwarded-For - if for some reason X-Real-Ip is blank and does not return response, get from X-Forwarded-For
		ipAddress = r.Header.Get("X-Forwarded-For")
	}

	if ipAddress == "" {
		// Remote Address - last resort (usually won't be reliable as this might be the last ip or if it is a naked http request to server ie no load balancer)
		ipAddress = r.RemoteAddr
	}

	return ipAddress
}
