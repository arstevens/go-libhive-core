package message

const (
	CapsuleType = "capsule"
)

func Decapsulate(msg *Message) ([]*Message, error) {
	layers := make([]*Message, 0)

	curMsg := msg
	mBody := curMsg.Body()
	layers = append(layers, curMsg)
	for curMsg.Header().Type() == "capsule" {
		nHeader := NewBufferedMessageHeader(mBody)
		curMsg = NewMessage(nHeader, mBody)
		mBody = curMsg.Body()
		layers = append(layers, curMsg)
	}
	return layers, nil
}
