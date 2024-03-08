package utils

import (
	"bufio"
	"os"
	"strconv"
	"zeus/models/drawBuffer"
	"zeus/models/row"
	"zeus/models/state"
)

func RefreshBuffer(buffer *drawBuffer.DrawBuffer, state *state.State) {
	buffer.Reset()
	for i := 0; i < state.TextLines; i++ {
		buffer.AppendText(state.Text[i].Text)
	}
	buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))
}

func OpenFile(filePath string, buffer *drawBuffer.DrawBuffer, state *state.State) {
	readFile, err := os.Open(filePath)

	defer readFile.Close()

	if err != nil {
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		r := row.NewRow(fileScanner.Text())
		state.AppendRow(&r)
	}
	state.UpdateLineno()
	RefreshBuffer(buffer, state)
}
