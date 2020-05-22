package versions

type DbbcMessage interface {
	UnmarshalBinary([]byte) error
}

var Messages = map[string]DbbcMessage{}

func Add(version string, message DbbcMessage) {
	Messages[version] = message
}
