package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"zeus/models/state"
	"zeus/models/row"
	"zeus/models/drawBuffer"
  "zeus/utils"
  "os"
)

// initial state of the editor
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
  // configure the window
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "Editor")
	defer rl.CloseWindow()
  
  // set the exit key of the window to null
	rl.SetExitKey(rl.KeyNull)

  // set all the colors in the editor
	backgroundColor := rl.NewColor(19, 16, 32, 255)
	textColor := rl.NewColor(243, 188, 230, 255)
	cursorColor := rl.NewColor(194, 157, 241, 255)
	linenoColor := rl.White
	linenoBackground := rl.NewColor(113, 98, 176, 255)
  
  // set the initial state of the window
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
  

  // initialize the draw buffer that handles drawing all the text to the screen
	var buffer *drawBuffer.DrawBuffer = drawBuffer.NewBuffer("","")

  // if a file name is passed then it is opened in the editor
  if len(os.Args) > 1 {
    filePath := os.Args[1]
    utils.OpenFile(filePath, buffer, state)

    // if there is no file name passed then it opens an empty editor
  } else {
    r := row.NewRow("")
    state.AppendRow(&r)
    state.UpdateLineno()
    utils.RefreshBuffer(buffer, state)
  }
	// set the cursor's horizontal position
	state.CursorX = state.LinenoOff
  
  // main loop that currently runs on an uncapped fps
	for !rl.WindowShouldClose() {

    // check if the right arrow key is pressed
		if (rl.IsKeyPressedRepeat(rl.KeyRight) || rl.IsKeyPressed(rl.KeyRight)) {

      // if the cursor is at the end of the current line and there are lines after then the cursor is pushed to the beginning  of the next line
			if state.CursorX == (int32(len(state.Text[state.GetCurrentRow()].Text)) * state.CursorHorizontalJump) + state.LinenoOff && state.IsCursorWithinText() {
				state.MoveCursorDown()
				state.CursorX = state.LinenoOff

        // if the cursor if within the current line text then move it to the next character
			} else if state.IsCursorWithinLine() {
				state.MoveCursorRight()
			}
		}

    // check if the left arrow key is pressed
		if (rl.IsKeyPressedRepeat(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyLeft)) {

      // if the cursor is at the beginning of the current line and there are lines before then the cursor is pushed to the end of the previous line
			if state.CursorX == state.LinenoOff && state.GetCurrentRow() != 0 {
				prev := state.GetCurrentRow()-1
				state.CursorX = (int32(len(state.Text[prev].Text)) * state.CursorHorizontalJump) + state.LinenoOff
				state.MoveCursorUp()

        // if the cursor is within the text of the current line then it is pushed to the previous character
			} else if state.CursorX != state.LinenoOff {
				state.MoveCursorLeft()
			}
		}

    // check if the up arrow key is pressed and push the cursor to the previous line if it exists 
		if (rl.IsKeyPressedRepeat(rl.KeyUp) || rl.IsKeyPressed(rl.KeyUp)) && state.CursorY != 0 {
			prev := state.GetCurrentRow() - 1
			l := len(state.Text[prev].Text)
			curr := state.GetCurrentCol()
			if int(curr) > l {
				state.CursorX = (int32(l) * state.CursorHorizontalJump) + state.LinenoOff
			}
			state.MoveCursorUp()
		}

    // check if the down arrow key is pressed and push the cursor down if is a next line
		if (rl.IsKeyPressedRepeat(rl.KeyDown) || rl.IsKeyPressed(rl.KeyDown)) && state.IsCursorWithinText() {
			next := state.GetCurrentRow() + 1
			l := len(state.Text[next].Text)
			curr := state.GetCurrentCol()
			if int(curr) > l {
				state.CursorX = (int32(l) * state.CursorHorizontalJump) + state.LinenoOff
			}
			state.MoveCursorDown()
		}

    // get the character pressed
		key := rl.GetCharPressed()

    // check if it is actually a character that was pressed
		for key > 0 {

      // check if the character pressed is a printable character and if it is then add it to the state text
			if key >= 32 && key <= 125 {
				colIndex := state.GetCurrentCol()
				text := state.Text[state.GetCurrentRow()].Text
				state.Text[state.GetCurrentRow()].Text = text[:colIndex] + string(rune(key)) + text[colIndex:]
				utils.RefreshBuffer(buffer, state)
        state.MoveCursorRight()
			}
      
      // get the keys they are registereed within the same frame
			key = rl.GetCharPressed()
		}

    // check if the backspace key is pressed
		if (rl.IsKeyPressedRepeat(rl.KeyBackspace) || rl.IsKeyPressed(rl.KeyBackspace)) {

      // if the cursor is at the beginning of the line then append it to the previous line if it exists
			if state.CursorX == state.LinenoOff && state.GetCurrentRow() != 0 {
				l := len(state.Text[state.GetCurrentRow()-1].Text)
				state.Text[state.GetCurrentRow()].Text = state.Text[state.GetCurrentRow()-1].Text + state.Text[state.GetCurrentRow()].Text
				copy(state.Text[state.GetCurrentRow()-1:], state.Text[state.GetCurrentRow():])
				state.Text = state.Text[:len(state.Text)-1]
				state.TextLines -= 1
        state.UpdateLineno()
				utils.RefreshBuffer(buffer, state)
        state.MoveCursorUp()
				state.CursorX = (int32(l) * state.CursorHorizontalJump) + state.LinenoOff

        // if the cursor is inbetween the current line then remove the prev character
			} else if state.CursorX != state.LinenoOff {
				colIndex := state.GetCurrentCol()
				text := state.Text[state.GetCurrentRow()].Text
				state.Text[state.GetCurrentRow()].Text = text[:colIndex-1] + text[colIndex:]
				utils.RefreshBuffer(buffer, state)
        state.MoveCursorLeft()
			}
		}

    // check if the enter key is pressed and append the current characters after the cursor to the next line
		if rl.IsKeyPressedRepeat(rl.KeyEnter) || rl.IsKeyPressed(rl.KeyEnter) {
			r := row.NewRow("")
			state.AppendRow(&r)
			copy(state.Text[state.GetCurrentRow()+1:], state.Text[state.GetCurrentRow():])
			temp := state.Text[state.GetCurrentRow()].Text[state.GetCurrentCol():]
			state.Text[state.GetCurrentRow()].Text = state.Text[state.GetCurrentRow()].Text[:state.GetCurrentCol()]
			state.Text[state.GetCurrentRow()+1].Text = temp
      state.UpdateLineno()
			utils.RefreshBuffer(buffer, state)
      state.MoveCursorDown()
			state.CursorX = state.LinenoOff
		}
		
    rl.BeginDrawing()

		rl.ClearBackground(state.BackgroundColor)
    
    // draws the cursor on to the screen
		rl.DrawRectangle(state.CursorX, state.CursorY, int32(state.TextCursorSize.X), int32(state.TextCursorSize.Y), state.CursorColor)
		
    // draws the line no column
		for i := 0; i < state.TextLines ; i++ {
			rl.DrawRectangle(state.TextX, state.TextY, int32(state.LinenoOff - int32(state.TextCursorSize.X)), int32(state.TextCursorSize.Y), state.LinenoBackground)
			state.TextX = 0
			state.TextY += state.CursorVerticalJump
		}
		
		state.TextX = 0
		state.TextY = 0

    // draws the line nos to the screen
		rl.DrawTextEx(state.Font, buffer.Lineno, rl.Vector2{X: float32(state.TextX), Y: float32(state.TextY)}, state.FontSize, state.Spacing, state.LinenoColor)

		state.TextX = state.LinenoOff
		state.TextY = 0
  
    // draws the text to the screen
		rl.DrawTextEx(state.Font, buffer.Text, rl.Vector2{X: float32(state.TextX), Y: float32(state.TextY)}, state.FontSize, state.Spacing, state.TextColor)
		state.TextX = 0
		state.TextY = 0

		rl.EndDrawing()
	}
}


