package protocol

const (
	ProtocolPrefix = "/x"
)

type ID string

func NewProtocolId(proto string) ID {
	if len(proto) > 0 && proto[0] == '/' {
		return ID(ProtocolPrefix + proto)
	}
	return ID(ProtocolPrefix + "/" + proto)
}

func (i ID) String() string { return string(i) }
