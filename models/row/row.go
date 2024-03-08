package row

type Row struct {
	Text string
}

func NewRow(text string) Row {
	return Row {
		Text: text,
	}
} 
