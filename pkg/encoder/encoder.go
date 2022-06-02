package encoder

// Encoder -
type Encoder interface {
	// Encode -
	Encode(interface{}) ([]byte, error)

	// Decode -
	Decode([]byte, interface{}) error

	// Type -
	Type() string
}
