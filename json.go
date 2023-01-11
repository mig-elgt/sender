package sender

import (
	"encoding/json"
	"net/http"

	"github.com/mig-elgt/sender/codes"
	"github.com/pkg/errors"
)

// jsonSender holds the reponse writer to set the http response data
// about the request.
type jsonSender struct {
	w          http.ResponseWriter
	statusCode int
	withErr    bool
	err        responseError
}

// NewJSON creates new instance of jsonSender struct.
func NewJSON(w http.ResponseWriter, code int) *jsonSender {
	return &jsonSender{
		w:          w,
		statusCode: code,
	}
}

// responseError describes the response structure for the errors.
type responseError struct {
	Status      int               `json:"status"`
	Error       string            `json:"error"`
	Description string            `json:"description,omitempty"`
	Fields      map[string]string `json:"fields,omitempty"`
}

// WithError sets a new response error.
func (js *jsonSender) WithError(code codes.Code, description string) *jsonSender {
	js.withErr = true
	js.err = responseError{
		Status:      js.statusCode,
		Error:       codes.ToString[code],
		Description: description,
	}
	return js
}

// WithFieldError sets a validation error as default with a field custom error.
func (js *jsonSender) WithFieldError(code codes.Code, field, value string) *jsonSender {
	js.WithError(code, "One or more fields raised validation errors.")
	js.err.Fields = map[string]string{field: value}
	return js
}

// WithFieldsError registers a set of custom errors and return a self pointer object.
func (js *jsonSender) WithFieldsError(code codes.Code, fields map[string]string) *jsonSender {
	js.WithError(code, "One or more fields raised validation errors.")
	js.err.Fields = fields
	return js
}

// Send sends an response error if the sender has one, otherwise
// response the content value as JSON object.
func (js *jsonSender) Send(content ...interface{}) error {
	js.w.Header().Set("Content-Type", "application/json")
	js.w.WriteHeader(js.statusCode)
	if js.withErr {
		type data struct {
			Error responseError `json:"error"`
		}
		if err := json.NewEncoder(js.w).Encode(&data{Error: js.err}); err != nil {
			return errors.Wrapf(err, "could not encode error response: %v", content)
		}
	}
	if len(content) > 0 {
		if err := json.NewEncoder(js.w).Encode(&content[0]); err != nil {
			return errors.Wrapf(err, "could not encode json response: %v", content)
		}
	}
	return nil
}
