# HTTP JSON Sender

Sender is a GO library for create a HTTP JSON response using http.ResponseWriter as writer output.

## Features

* Send a JSON response error with a simple description error.
* Send a JSON response error with a custom field error.
* Send a JSON response error with a list of custom fields errors.
* Send a JSON response object.

# Installation

Standard `go get`:

```
$ go get github.com/mig-elgt/sender
```

## Usage & Example

### Handle to send a simple error description

```go
func WithErrorHandle(w http.ResponseWriter, r *http.Request) {
	sender.
		NewJSON(w, http.StatusInternalServerError).
		WithError(codes.Internal, "Something went wrong..").
		Send()
}
// Output
// {
//     "error": {
//         "status": 500,
//         "error": "INTERNAL",
//         "description": "Something went wrong..."
//     }
// }
```

### Handle to send a simple error with a field value.

```go
func WithFieldErrorHandle(w http.ResponseWriter, r *http.Request) {
	sender.
		NewJSON(w, http.StatusBadRequest).
		WithFieldError(
			codes.InvalidArgument,
			"user_id",
			"User ID is required",
		).
		Send()
}

// Output
// {
//     "error": {
//         "status": 500,
//         "error": "INTERNAL",
//         "description": "Something went wrong..."
//     }
// }
```

### Handle to send an error with a set of fields inputs.

``` go
func WithFieldsErrorHandle(w http.ResponseWriter, r *http.Request) {
	sender.
		NewJSON(w, http.StatusBadRequest).
		WithFieldsError(
			codes.InvalidArgument,
			map[string]string{
				"user_id": "User ID is required",
				"email":   "Email has invalid format",
			},
		).
		Send()
}

// Output
// {
//     "error": {
//         "status": 400,
//         "error": "INVALID_ARGUMENT",
//         "description": "One or more fields raised validation errors.",
//         "fields": {
//             "email": "Email has invalid format",
//             "user_id": "User ID is required"
//         }
//     }
// }
```
