package object

// Key is used to identify unique values in a set.
type Key struct {
	// Type of the object being referenced.
	Type Type
	// FltValue is used as the key for floats.
	FltValue float64
	// IntValue is used as the key for integers.
	IntValue int64
	// StrValue is used as the key for strings.
	StrValue string
}
