package editor

import "errors"

type Buffer struct {
	lines [][]rune
}

func NewBuffer(lines []string) *Buffer {
	buf := &Buffer{}
	if len(lines) == 0 {
		buf.lines = [][]rune{[]rune{}}
		return buf
	}
	buf.lines = make([][]rune, len(lines))
	for i, line := range lines {
		buf.lines[i] = []rune(line)
	}
	return buf
}

func (b *Buffer) LineCount() int {
	return len(b.lines)
}

func (b *Buffer) Line(index int) ([]rune, error) {
	if index < 0 || index >= len(b.lines) {
		return nil, errors.New("line out of range")
	}
	return b.lines[index], nil
}

func (b *Buffer) InsertRune(row int, col int, value rune) error {
	line, err := b.Line(row)
	if err != nil {
		return err
	}
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}
	line = append(line[:col], append([]rune{value}, line[col:]...)...)
	b.lines[row] = line
	return nil
}

func (b *Buffer) DeleteRune(row int, col int) error {
	line, err := b.Line(row)
	if err != nil {
		return err
	}
	if col < 0 || col >= len(line) {
		return errors.New("column out of range")
	}
	line = append(line[:col], line[col+1:]...)
	b.lines[row] = line
	return nil
}

func (b *Buffer) SplitLine(row int, col int) error {
	line, err := b.Line(row)
	if err != nil {
		return err
	}
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}
	left := append([]rune{}, line[:col]...)
	right := append([]rune{}, line[col:]...)
	b.lines[row] = left
	b.lines = append(b.lines[:row+1], append([][]rune{right}, b.lines[row+1:]...)...)
	return nil
}

func (b *Buffer) MergeLine(row int) error {
	if row < 0 || row+1 >= len(b.lines) {
		return errors.New("line out of range")
	}
	b.lines[row] = append(b.lines[row], b.lines[row+1]...)
	b.lines = append(b.lines[:row+1], b.lines[row+2:]...)
	return nil
}

func (b *Buffer) LinesAsStrings() []string {
	result := make([]string, len(b.lines))
	for i, line := range b.lines {
		result[i] = string(line)
	}
	return result
}
