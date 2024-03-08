package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"zeus/models/state"
	"zeus/models/row"
	"zeus/models/drawBuffer"
	"strconv"
)

const (
	screenWidth = 1920
	screenHeight = 1080
	spacing = 2
	cursorX = 0
	cursorY = 0
	textX = 0
	textY = 0
	fontSize = 17.0
	fontPath = "font/JetBrainsMono-Medium.ttf"
)



func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "Editor")
	defer rl.CloseWindow()

	rl.SetExitKey(rl.KeyNull)

	backgroundColor := rl.NewColor(19, 16, 32, 255)
	textColor := rl.NewColor(243, 188, 230, 255)
	cursorColor := rl.NewColor(194, 157, 241, 255)
	linenoColor := rl.White
	linenoBackground := rl.NewColor(113, 98, 176, 255)

	state := state.InitState(
		fontPath,
		cursorX, 
		cursorY,
		textX,
		textY,
		textColor, 
		cursorColor,
		backgroundColor,
		linenoColor,
		linenoBackground,
		screenHeight,
		screenWidth,
		spacing,
		fontSize,
	)
	state.ScreenWidth = int32(rl.GetScreenWidth())
	sampleText := []string{"Bruv", "Whats up?", "Nin Amman", "Nin Amman", "Nin Amman", "Nin Amman", "Nin Amman", "Nin Amman", "Nin Amman", "Nin Amman", "Nin Amman"}

	for _ , s := range sampleText {
		r := row.NewRow(s)
		state.AppendRow(&r)
	}

	state.LinenoOff = int32((state.CursorHorizontalJump * int32(len(strconv.Itoa(state.TextLines)) + 1)))
	state.CursorX = state.LinenoOff
	state.VisibleCols = state.GetVisibleCols()
	state.VisibleCols -= 10
	state.VisibleRows = state.GetVisibleRows()

	var buffer *drawBuffer.DrawBuffer = drawBuffer.NewBuffer("","")

	for !rl.WindowShouldClose() {
		if (rl.IsKeyPressedRepeat(rl.KeyRight) || rl.IsKeyPressed(rl.KeyRight)) {
			if state.CursorX == (int32(len(state.Text[state.GetCurrentRow()].Text)) * state.CursorHorizontalJump) + state.LinenoOff - (state.ColOff * state.CursorHorizontalJump) && state.IsCursorWithinText() {
				state.MoveCursorDown()
				state.CursorX = state.LinenoOff
				state.ColOff = 0
				buffer.ColOff = 0
			} else if state.IsCursorWithinLine() {
				state.MoveCursorRight()
				curr := state.GetCurrentCol()
				if curr > state.VisibleCols {
					state.ColOff += (curr - state.VisibleCols)
					buffer.ColOff = state.ColOff
					state.MoveCursorLeft()
				}
			}
		}
		if (rl.IsKeyPressedRepeat(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyLeft)) {
			if state.CursorX == state.LinenoOff && state.GetCurrentRow() != 0 && state.ColOff == 0 {
				prev := state.GetCurrentRow()-1
				state.CursorX = (int32(len(state.Text[prev].Text)) * state.CursorHorizontalJump) + state.LinenoOff
				state.MoveCursorUp()
				state.ColOff = 0
				buffer.ColOff = 0
				curr := state.GetCurrentCol()
				if curr > state.VisibleCols {
					state.ColOff += (curr - state.VisibleCols + 1)
					buffer.ColOff = state.ColOff
					state.MoveCursorLeft()
				}
			} else if state.CursorX != state.LinenoOff {
				state.MoveCursorLeft()
			} else if state.IsCursorWithinLine() && state.CursorX == state.LinenoOff {
				if state.ColOff != 0 {
					state.ColOff -= 1
					buffer.ColOff = state.ColOff
					state.MoveCursorLeft()
				}
			}
		}
		if (rl.IsKeyPressedRepeat(rl.KeyUp) || rl.IsKeyPressed(rl.KeyUp)) && state.CursorY != 0 {
			prev := state.GetCurrentRow() - 1
			l := len(state.Text[prev].Text)
			curr := state.GetCurrentCol()
			if int(curr) > l {
				state.ColOff = 0
				buffer.ColOff = 0
				state.CursorX = (int32(l) * state.CursorHorizontalJump) + state.LinenoOff - (state.ColOff * state.CursorHorizontalJump)
			}
			state.MoveCursorUp()
		}
		if (rl.IsKeyPressedRepeat(rl.KeyDown) || rl.IsKeyPressed(rl.KeyDown)) && state.IsCursorWithinText() {
			next := state.GetCurrentRow() + 1
			l := len(state.Text[next].Text)
			curr := state.GetCurrentCol()
			if int(curr) > l {
				state.CursorX = (int32(l) * state.CursorHorizontalJump) + state.LinenoOff - (state.ColOff * state.CursorHorizontalJump)
			}
			state.MoveCursorDown()
		}
		key := rl.GetCharPressed()
		for key > 0 {
			if key >= 32 && key <= 125 {
				colIndex := state.GetCurrentCol()
				text := state.Text[state.GetCurrentRow()].Text
				state.Text[state.GetCurrentRow()].Text = text[:colIndex] + string(rune(key)) + text[colIndex:]
				buffer.Reset()
				for _ , r := range state.Text {
					buffer.AppendText(r.Text, state.VisibleCols)
				}
				buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))
				state.MoveCursorRight()
				curr := state.GetCurrentCol()
				if curr > state.VisibleCols {
					state.ColOff += 1
					buffer.ColOff = state.ColOff
					state.MoveCursorLeft()
				}
			}
			key = rl.GetCharPressed()
		}
		if (rl.IsKeyPressedRepeat(rl.KeyBackspace) || rl.IsKeyPressed(rl.KeyBackspace)) {
			if state.CursorX == state.LinenoOff && state.GetCurrentRow() != 0 {
				l := len(state.Text[state.GetCurrentRow()-1].Text)
				state.Text[state.GetCurrentRow()].Text = state.Text[state.GetCurrentRow()-1].Text + state.Text[state.GetCurrentRow()].Text
				copy(state.Text[state.GetCurrentRow()-1:], state.Text[state.GetCurrentRow():])
				state.Text = state.Text[:len(state.Text)-1]
				state.TextLines -= 1
				buffer.Reset()
				for _ , r := range state.Text {
					buffer.AppendText(r.Text, state.VisibleCols)
				}
				buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))
				state.MoveCursorUp()
				state.CursorX = (int32(l) * state.CursorHorizontalJump) + state.LinenoOff
			} else {
				colIndex := state.GetCurrentCol()
				text := state.Text[state.GetCurrentRow()].Text
				state.Text[state.GetCurrentRow()].Text = text[:colIndex-1] + text[colIndex:]
				buffer.Reset()
				for _ , r := range state.Text {
					buffer.AppendText(r.Text, state.VisibleCols)
				}
				buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))
				state.MoveCursorLeft()
			}
		}
		if rl.IsKeyPressedRepeat(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyEnter) {
			r := row.NewRow("")
			state.AppendRow(&r)
			copy(state.Text[state.GetCurrentRow()+1:], state.Text[state.GetCurrentRow():])
			temp := state.Text[state.GetCurrentRow()].Text[state.GetCurrentCol():]
			state.Text[state.GetCurrentRow()].Text = state.Text[state.GetCurrentRow()].Text[:state.GetCurrentCol()]
			state.Text[state.GetCurrentRow()+1].Text = temp
			buffer.Reset()
			for _ , r := range state.Text {
				buffer.AppendText(r.Text, state.VisibleCols)
			}
			buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))
			state.MoveCursorDown()
			state.CursorX = state.LinenoOff

		}

		buffer.Reset()
		for i := state.RowOff; i < int32(state.TextLines); i++ {
			buffer.AppendText(state.Text[i].Text, state.VisibleCols)
		}
		buffer.AppendLineno(state.TextLines, len(strconv.Itoa(state.TextLines)))

		rl.BeginDrawing()

		rl.ClearBackground(state.BackgroundColor)

		rl.DrawRectangle(state.CursorX, state.CursorY, int32(state.TextCursorSize.X), int32(state.TextCursorSize.Y), state.CursorColor)
		
		for i := 0; i < state.TextLines ; i++ {
			rl.DrawRectangle(state.TextX, state.TextY, int32(state.LinenoOff - int32(state.TextCursorSize.X)), int32(state.TextCursorSize.Y), state.LinenoBackground)
			state.TextX = 0
			state.TextY += state.CursorVerticalJump
		}
		
		state.TextX = 0
		state.TextY = 0

		rl.DrawTextEx(state.Font, buffer.Lineno, rl.Vector2{X: float32(state.TextX), Y: float32(state.TextY)}, state.FontSize, state.Spacing, state.LinenoColor)

		state.TextX = state.LinenoOff
		state.TextY = 0

		rl.DrawTextEx(state.Font, buffer.Text, rl.Vector2{X: float32(state.TextX), Y: float32(state.TextY)}, state.FontSize, state.Spacing, state.TextColor)
		state.TextX = 0
		state.TextY = 0
		rl.EndDrawing()
	}
}


