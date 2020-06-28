package errmsg

type BadMsg struct {
	s string
}

func NewBadMsg(text string) error {
	return &BadMsg{text}
}

func (e *BadMsg) Error() string {
	return e.s
}
