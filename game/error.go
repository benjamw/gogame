package game

// KLUDGE: this should probably be in the golibs instead of the game module,
// because golibs references this, and it's more abstract than the game module.

// Error is an error implementer that also has a code value
type Error interface {
	// Error is the same function that the error interface requires
	// so this interface is an extension of the error interface
	Error() string

	// Code returns the HTTP Status Code for the error
	Code() int
}
