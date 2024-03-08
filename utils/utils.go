package utils

import (
  "zeus/models/drawBuffer"
  "zeus/models/state"
  "strconv"
)

func RefreshBuffer(buffer *drawBuffer.DrawBuffer, state *state.State) {
  buffer.Reset()
	for i := 0; i < state.TextLines; i++ {
		buffer.AppendText(state.Text[i].Text)
	}
	buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))
}


