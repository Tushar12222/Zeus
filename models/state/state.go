package state

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"zeus/models/row"
)

type State struct {
	Font rl.Font
	TextCursorSize rl.Vector2
	CursorHorizontalJump int32
	CursorVerticalJump int32
	CursorX int32
	CursorY int32
	TextX int32
	TextY int32
	ScreenHeight int32
	ScreenWidth int32
	Text []row.Row
	Modifications int
	Spacing float32
	FontSize float32
	TextColor rl.Color
	CursorColor rl.Color
	BackgroundColor rl.Color
	LinenoColor rl.Color
	LinenoBackground rl.Color
	TextLines int
	LinenoOff int32
	ColOff int32
	RowOff int32
	VisibleCols int32
	VisibleRows int32
}

func InitState(fontPath string, textX, textY, cursorX, cursorY int32, textColor, cursorColor, backgroundColor, linenoColor, linenoBackground rl.Color, screenHeight, screenWidth int32, spacing , fontSize float32) *State {
	font := rl.LoadFontEx(fontPath, int32(fontSize), nil)
	textCursorSize := rl.Vector2Add(rl.MeasureTextEx(font, "A", fontSize, spacing) , rl.Vector2{X: 1, Y: -2})
	return &State {
		Font: font,
		TextCursorSize: textCursorSize,
		CursorHorizontalJump: int32(textCursorSize.X + (spacing / 2)),
		CursorVerticalJump: int32(textCursorSize.Y),
		CursorX: cursorX,
		CursorY: cursorY,
		TextX: textX,
		TextY: textY,
		ScreenHeight: screenHeight,
		ScreenWidth: screenWidth,
		Text: []row.Row{},
		Modifications: 0,
		Spacing: spacing,
		FontSize: fontSize,
		TextColor: textColor,
		CursorColor: cursorColor,
		BackgroundColor: backgroundColor,
		LinenoColor: linenoColor,
		LinenoBackground: linenoBackground,
		TextLines: 0,
		LinenoOff: 0,
		ColOff: 0,
		RowOff: 0,
		VisibleCols: 0,
		VisibleRows: 0,
	}
}

func (state *State) AppendRow(row *row.Row) {
	state.Text = append(state.Text, *row)
	state.TextLines += 1
}

func (state *State) MoveCursorRight() {
	state.CursorX += state.CursorHorizontalJump
}

func (state *State) MoveCursorLeft() {
	state.CursorX -= state.CursorHorizontalJump
}

func (state *State) MoveCursorUp() {
	state.CursorY -= state.CursorVerticalJump
}

func (state *State) MoveCursorDown() {
	state.CursorY += state.CursorVerticalJump
}

func (state *State) GetCurrentCol() int32 {
	result := (state.CursorX - state.LinenoOff) / state.CursorHorizontalJump
	if result < 0 {
		return 0
	}
	return result
}

func (state *State) GetCurrentRow() int32 {
	result := (state.CursorY/ state.CursorVerticalJump)
	if result < 0 {
		return 0
	}
	return result
}

func (state *State) IsCursorWithinLine() bool {
	return int(state.GetCurrentCol() + state.ColOff) < len(state.Text[state.GetCurrentRow()].Text)
}

func (state *State) IsCursorWithinText() bool {
	return int(state.GetCurrentRow()) < state.TextLines-1
}

func (state *State) GetVisibleCols() int32 {
	return (state.ScreenWidth - state.LinenoOff) / state.CursorHorizontalJump
}

func (state *State) GetVisibleRows() int32 {
	return state.ScreenHeight / state.CursorVerticalJump
}
