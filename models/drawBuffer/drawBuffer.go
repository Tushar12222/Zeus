package drawBuffer

import (
	"strconv"
)

type DrawBuffer struct {
	Text string
	Lineno string
	ColOff int32
	RowOff int32
}

func NewBuffer(lineno, text string) *DrawBuffer {
	return &DrawBuffer{
		Text: text,
		Lineno: lineno,
		ColOff: 0,
		RowOff: 0,
	}
}

func (buffer *DrawBuffer) AppendText(text string, visibleCols int32) {
	l := len(text) - int(buffer.ColOff)
	if l <= 0 {
		buffer.Text = ""
	} else {
		buffer.Text += text[buffer.ColOff: int(buffer.ColOff)+l]
	}
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

