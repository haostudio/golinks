package encoding

// Binary defines a binary encoding/decoding interface.
type Binary interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}
