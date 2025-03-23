package domain

type ChatState string

const (
	StateEnterEmail ChatState = "ENTER_EMAIL"
	StateEnterCode  ChatState = "ENTER_CODE"
)

func (chatState ChatState) MarshalBinary() (data []byte, err error) {
	data = []byte(chatState)
	return
}

func (chatState ChatState) Valid() bool {
	if chatState == StateEnterEmail {
		return true
	}

	if chatState == StateEnterCode {
		return true
	}

	return false
}

func (chatState ChatState) String() string {
	return string(chatState)
}
