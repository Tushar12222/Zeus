package drawBuffer

import (
	"strconv"
)

type DrawBuffer struct {
	Text string
	Lineno string
}

func NewBuffer(lineno, text string) *DrawBuffer {
	return &DrawBuffer{
		Text: text,
		Lineno: lineno,
	}
}

func (buffer *DrawBuffer) AppendText(text string) {
  buffer.Text += text
	buffer.Text += "\n"
}

func (buffer *DrawBuffer) AppendLineno(totalLines int, colLen int) {
	for i := 1; i <= totalLines; i++ {
		s := strconv.Itoa(i)
		padding := colLen - len(s)
		for ; padding > 0 ; padding-- {
			buffer.Lineno += " "
		}
		buffer.Lineno += s
		buffer.Lineno += "\n"
	}
}

func (buffer *DrawBuffer) Reset() {
	buffer.Text = ""
	buffer.Lineno = ""
}

