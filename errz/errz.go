// Package errz defines a FriendlyError interface for errors that have a human
// friendly message in addition to the default error message.
package errz

// FriendlyError is an interface for errors that have a human friendly message
// in addition to a the lower level default error message.
type FriendlyError interface {
	Error() string
	FriendlyErrorMessage() string
}
