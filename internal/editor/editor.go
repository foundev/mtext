package editor

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"mtext/internal/terminal"
)

const (
	statusMessageDuration = 5 * time.Second
)

type Cursor struct {
	Row int
	Col int
}

type Editor struct {
	term       terminal.Terminal
	keyReader  *KeyReader
	output     io.Writer
	buffer     *Buffer
	cursor     Cursor
	rowOffset  int
	colOffset  int
	screenRows int
	screenCols int
	filename   string
	dirty      bool
	status     string
	statusTime time.Time
}

func New(term terminal.Terminal, input io.Reader, output io.Writer) (*Editor, error) {
	rows, cols, err := term.Size()
	if err != nil {
		return nil, err
	}
	if rows < 2 {
		return nil, errors.New("terminal too small")
	}
	return &Editor{
		term:       term,
		keyReader:  NewKeyReader(input),
		output:     output,
		buffer:     NewBuffer(nil),
		screenRows: rows - 2,
		screenCols: cols,
	}, nil
}

func (e *Editor) OpenFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	e.buffer = NewBuffer(lines)
	e.filename = path
	e.dirty = false
	return nil
}

func (e *Editor) Save() error {
	if e.filename == "" {
		return errors.New("no filename set")
	}
	content := strings.Join(e.buffer.LinesAsStrings(), "\n")
	if err := os.WriteFile(e.filename, []byte(content), 0o644); err != nil {
		return err
	}
	e.dirty = false
	e.SetStatusMessage("%d bytes written", len(content))
	return nil
}

func (e *Editor) Run() error {
	if err := e.term.EnableRawMode(); err != nil {
		return err
	}
	defer e.term.DisableRawMode()

	for {
		if err := e.RefreshScreen(); err != nil {
			return err
		}
		key, err := e.keyReader.ReadKey()
		if err != nil {
			return err
		}
		if err := e.ProcessKey(key); err != nil {
			return err
		}
	}
}

func (e *Editor) ProcessKey(key Key) error {
	switch key.Kind {
	case KeyCtrlQ:
		return io.EOF
	case KeyCtrlS:
		if err := e.Save(); err != nil {
			e.SetStatusMessage("save failed: %v", err)
		}
	case KeyArrowUp:
		e.moveCursor(-1, 0)
	case KeyArrowDown:
		e.moveCursor(1, 0)
	case KeyArrowLeft:
		e.moveCursor(0, -1)
	case KeyArrowRight:
		e.moveCursor(0, 1)
	case KeyHome:
		e.cursor.Col = 0
	case KeyEnd:
		lineLen := e.currentLineLength()
		e.cursor.Col = lineLen
	case KeyPageUp:
		e.cursor.Row -= e.screenRows
		if e.cursor.Row < 0 {
			e.cursor.Row = 0
		}
	case KeyPageDown:
		e.cursor.Row += e.screenRows
		if e.cursor.Row >= e.buffer.LineCount() {
			e.cursor.Row = e.buffer.LineCount() - 1
		}
	case KeyBackspace:
		if err := e.deleteChar(); err != nil {
			return err
		}
	case KeyDelete:
		e.moveCursor(0, 1)
		if err := e.deleteChar(); err != nil {
			return err
		}
	case KeyEnter:
		if err := e.insertNewline(); err != nil {
			return err
		}
	case KeyRune:
		if err := e.insertRune(key.Rune); err != nil {
			return err
		}
	}
	return nil
}

func (e *Editor) insertRune(value rune) error {
	if err := e.buffer.InsertRune(e.cursor.Row, e.cursor.Col, value); err != nil {
		return err
	}
	e.cursor.Col++
	e.dirty = true
	return nil
}

func (e *Editor) insertNewline() error {
	if err := e.buffer.SplitLine(e.cursor.Row, e.cursor.Col); err != nil {
		return err
	}
	e.cursor.Row++
	e.cursor.Col = 0
	e.dirty = true
	return nil
}

func (e *Editor) deleteChar() error {
	if e.cursor.Row == 0 && e.cursor.Col == 0 {
		return nil
	}
	if e.cursor.Col > 0 {
		if err := e.buffer.DeleteRune(e.cursor.Row, e.cursor.Col-1); err != nil {
			return err
		}
		e.cursor.Col--
		e.dirty = true
		return nil
	}
	prevLen := e.currentLineLengthAt(e.cursor.Row - 1)
	if err := e.buffer.MergeLine(e.cursor.Row - 1); err != nil {
		return err
	}
	e.cursor.Row--
	e.cursor.Col = prevLen
	e.dirty = true
	return nil
}

func (e *Editor) moveCursor(dRow int, dCol int) {
	e.cursor.Row += dRow
	e.cursor.Col += dCol

	if e.cursor.Row < 0 {
		e.cursor.Row = 0
	}
	if e.cursor.Row >= e.buffer.LineCount() {
		e.cursor.Row = e.buffer.LineCount() - 1
	}

	lineLen := e.currentLineLength()
	if e.cursor.Col < 0 {
		if e.cursor.Row > 0 {
			e.cursor.Row--
			e.cursor.Col = e.currentLineLength()
		} else {
			e.cursor.Col = 0
		}
	}
	if e.cursor.Col > lineLen {
		e.cursor.Col = lineLen
	}
}

func (e *Editor) currentLineLength() int {
	return e.currentLineLengthAt(e.cursor.Row)
}

func (e *Editor) currentLineLengthAt(row int) int {
	line, err := e.buffer.Line(row)
	if err != nil {
		return 0
	}
	return len(line)
}

func (e *Editor) RefreshScreen() error {
	e.scroll()
	var builder strings.Builder
	builder.WriteString("\x1b[?25l")
	builder.WriteString("\x1b[H")
	e.drawRows(&builder)
	e.drawStatusBar(&builder)
	e.drawMessageBar(&builder)
	cursorRow := e.cursor.Row - e.rowOffset + 1
	cursorCol := e.cursor.Col - e.colOffset + 1
	builder.WriteString(fmt.Sprintf("\x1b[%d;%dH", cursorRow, cursorCol))
	builder.WriteString("\x1b[?25h")
	_, err := io.WriteString(e.output, builder.String())
	return err
}

func (e *Editor) drawRows(builder *strings.Builder) {
	for row := 0; row < e.screenRows; row++ {
		fileRow := row + e.rowOffset
		if fileRow >= e.buffer.LineCount() {
			if e.buffer.LineCount() == 1 && row == e.screenRows/3 {
				title := "mtext -- version 0.0.1"
				if len(title) > e.screenCols {
					title = title[:e.screenCols]
				}
				padding := (e.screenCols - len(title)) / 2
				builder.WriteString("~")
				for i := 0; i < padding-1; i++ {
					builder.WriteByte(' ')
				}
				builder.WriteString(title)
			} else {
				builder.WriteString("~")
			}
			builder.WriteString("\x1b[K")
			builder.WriteString("\r\n")
			continue
		}
		line, _ := e.buffer.Line(fileRow)
		rendered := string(line)
		if len(rendered) > e.colOffset {
			rendered = rendered[e.colOffset:]
		} else {
			rendered = ""
		}
		if len(rendered) > e.screenCols {
			rendered = rendered[:e.screenCols]
		}
		builder.WriteString(rendered)
		builder.WriteString("\x1b[K")
		builder.WriteString("\r\n")
	}
}

func (e *Editor) drawStatusBar(builder *strings.Builder) {
	builder.WriteString("\x1b[7m")
	name := e.filename
	if name == "" {
		name = "[No Name]"
	}
	status := fmt.Sprintf("%.20s - %d lines", name, e.buffer.LineCount())
	if e.dirty {
		status += " (modified)"
	}
	rstatus := fmt.Sprintf("%d/%d", e.cursor.Row+1, e.buffer.LineCount())
	if len(status) > e.screenCols {
		status = status[:e.screenCols]
	}
	builder.WriteString(status)
	for len(status)+len(rstatus) < e.screenCols {
		builder.WriteByte(' ')
		status += " "
	}
	builder.WriteString(rstatus)
	builder.WriteString("\x1b[m")
	builder.WriteString("\r\n")
}

func (e *Editor) drawMessageBar(builder *strings.Builder) {
	builder.WriteString("\x1b[K")
	if e.status == "" {
		builder.WriteString("\r\n")
		return
	}
	if time.Since(e.statusTime) > statusMessageDuration {
		e.status = ""
		builder.WriteString("\r\n")
		return
	}
	message := e.status
	if len(message) > e.screenCols {
		message = message[:e.screenCols]
	}
	builder.WriteString(message)
	builder.WriteString("\r\n")
}

func (e *Editor) scroll() {
	if e.cursor.Row < e.rowOffset {
		e.rowOffset = e.cursor.Row
	}
	if e.cursor.Row >= e.rowOffset+e.screenRows {
		e.rowOffset = e.cursor.Row - e.screenRows + 1
	}
	if e.cursor.Col < e.colOffset {
		e.colOffset = e.cursor.Col
	}
	if e.cursor.Col >= e.colOffset+e.screenCols {
		e.colOffset = e.cursor.Col - e.screenCols + 1
	}
}

func (e *Editor) SetStatusMessage(format string, args ...any) {
	e.status = fmt.Sprintf(format, args...)
	e.statusTime = time.Now()
}
