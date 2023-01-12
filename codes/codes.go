package codes

type Code uint32

const (
	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Code = 1
	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal Code = 2
	// InvalidArgument indicates client specified an invalid argument.
	InvalidArgument Code = 3

	NotFound Code = 4

	NotAuthorized = 5
)

var ToString = map[Code]string{
	AlreadyExists:   "ALREADY_EXISTS",
	Internal:        "INTERNAL",
	InvalidArgument: "INVALID_ARGUMENT",
	NotFound:        "NOT_FOUND",
	NotAuthorized:   "NOT_AUTHORIZED",
}
