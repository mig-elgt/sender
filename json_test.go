package sender

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mig-elgt/sender/codes"
)

func TestJsonSender_NewJSON(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		code int
	}
	testCases := []struct {
		name string
		args args
		want *jsonSender
	}{
		{
			name: "base-case",
			args: args{
				w:    httptest.NewRecorder(),
				code: http.StatusOK,
			},
			want: &jsonSender{
				w:          httptest.NewRecorder(),
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			js := NewJSON(tc.args.w, tc.args.code)
			if !reflect.DeepEqual(js, tc.want) {
				t.Errorf("NewJSON(w,code) got %v; want %v", js, tc.want)
			}
		})
	}
}

func TestJsonSender_WithError(t *testing.T) {
	type args struct {
		w       *httptest.ResponseRecorder
		code    int
		errCode codes.Code
		errDesc string
	}
	testCases := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "already exists",
			args: args{
				w:       httptest.NewRecorder(),
				code:    http.StatusConflict,
				errCode: codes.AlreadyExists,
				errDesc: "Email already exists",
			},
			want: []byte("{\"error\":{\"status\":409,\"error\":\"ALREADY_EXISTS\",\"description\":\"Email already exists\"}}\n"),
		},
		{
			name: "internal server error",
			args: args{
				w:       httptest.NewRecorder(),
				code:    http.StatusInternalServerError,
				errCode: codes.Internal,
				errDesc: "Something went wrong...",
			},
			want: []byte("{\"error\":{\"status\":500,\"error\":\"INTERNAL\",\"description\":\"Something went wrong...\"}}\n"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			NewJSON(tc.args.w, tc.args.code).WithError(tc.args.errCode, tc.args.errDesc).Send()
			if got, want := tc.args.w.Body.Bytes(), tc.want; !reflect.DeepEqual(got, want) {
				t.Errorf("WithError(code) got \n%+v; want \n%+v", string(got), string(want))
			}
		})
	}
}

func TestJsonSender_WithFieldError(t *testing.T) {
	type args struct {
		w          *httptest.ResponseRecorder
		code       int
		errCode    codes.Code
		fieldName  string
		fieldValue string
	}
	testCases := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "email invalid",
			args: args{
				w:          httptest.NewRecorder(),
				code:       http.StatusUnprocessableEntity,
				errCode:    codes.InvalidArgument,
				fieldName:  "email",
				fieldValue: "should be a string",
			},
			want: []byte("{\"error\":{\"status\":422,\"error\":\"INVALID_ARGUMENT\",\"description\":\"One or more fields raised validation errors.\",\"fields\":{\"email\":\"should be a string\"}}}\n"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			NewJSON(tc.args.w, tc.args.code).WithFieldError(tc.args.errCode, tc.args.fieldName, tc.args.fieldValue).Send()
			if got, want := tc.args.w.Body.Bytes(), tc.want; !reflect.DeepEqual(got, want) {
				t.Errorf("WithFieldError(code) got \n%+v; want \n%+v", string(got), string(want))
			}
		})
	}
}

func TestJsonSender_WithFieldsError(t *testing.T) {
	type args struct {
		w       *httptest.ResponseRecorder
		code    int
		errCode codes.Code
		fields  map[string]string
	}
	cases := map[string]struct {
		args args
		want []byte
	}{
		"fields missing": {
			args: args{
				w:       httptest.NewRecorder(),
				code:    http.StatusUnprocessableEntity,
				errCode: codes.InvalidArgument,
				fields:  map[string]string{"username": "The user name is blank.", "email": "The email is required."},
			},
			want: []byte("{\"error\":{\"status\":422,\"error\":\"INVALID_ARGUMENT\",\"description\":\"One or more fields raised validation errors.\",\"fields\":{\"email\":\"The email is required.\",\"username\":\"The user name is blank.\"}}}\n"),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			NewJSON(tc.args.w, tc.args.code).WithFieldsError(tc.args.errCode, tc.args.fields).Send()
			if got, want := tc.args.w.Body.Bytes(), tc.want; !reflect.DeepEqual(got, want) {
				t.Errorf("WithFieldsError(code, fields) got \n%v; want \n%v", string(got), string(want))
			}
		})
	}
}

// func TestJsonSender_Send(t *testing.T) {
// 	type newAccount struct {
// 		Email     string `json:"email"`
// 		Activated bool   `json:"activated"`
// 	}
// 	type args struct {
// 		w    *httptest.ResponseRecorder
// 		code int
// 		reg  newAccount
// 	}
// 	testCases := []struct {
// 		name string
// 		args args
// 		want []byte
// 	}{
// 		{
// 			name: "base-case",
// 			args: args{
// 				w:    httptest.NewRecorder(),
// 				code: http.StatusOK,
// 				reg:  newAccount{"miguel@gmail.com", false},
// 			},
// 			want: []byte(fmt.Sprintf("{\"email\":\"miguel@gmail.com\",\"activated\":false}\n")),
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			NewJSON(tc.args.w, tc.args.code).Send(tc.args.reg)
// 			if got, want := tc.args.w.Body.Bytes(), tc.want; !reflect.DeepEqual(got, want) {
// 				t.Errorf("WithError(code) got \n%+v; want \n%+v", string(got), string(want))
// 			}
// 		})
// 	}
// }

// func createWriter(content string) http.ResponseWriter {
// 	w := httptest.NewRecorder()
// 	w.Body.WriteString(content)
// 	return w
// }
