package message

import "fmt"

func Decapsulate(msg *Message) ([]*Message, error) {
	layers := make([]*Message, 0)

	curMsg := msg
	mBody := curMsg.Body()
	layers = append(layers, curMsg)
	for curMsg.Header().IsCapsule() {
		fmt.Println(curMsg.Body().Name())
		curMsg, err := ReadMessage(mBody)
		if err != nil {
			return layers, err
		}
		mBody = curMsg.Body()
		layers = append(layers, curMsg)
	}
	return layers, nil
}
