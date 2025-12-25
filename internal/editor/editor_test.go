package editor

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type fakeTerminal struct {
	rows int
	cols int
}

func (f fakeTerminal) EnableRawMode() error    { return nil }
func (f fakeTerminal) DisableRawMode() error   { return nil }
func (f fakeTerminal) Size() (int, int, error) { return f.rows, f.cols, nil }

func newTestEditor(t *testing.T) *Editor {
	t.Helper()
	ed, err := New(fakeTerminal{rows: 10, cols: 40}, bytes.NewBuffer(nil), io.Discard)
	if err != nil {
		t.Fatalf("new editor: %v", err)
	}
	return ed
}

func TestEditorInsertAndDelete(t *testing.T) {
	ed := newTestEditor(t)
	if err := ed.ProcessKey(RuneKey('a')); err != nil {
		t.Fatalf("insert rune: %v", err)
	}
	line, _ := ed.buffer.Line(0)
	if string(line) != "a" {
		t.Fatalf("expected buffer to contain 'a', got %q", string(line))
	}
	if err := ed.ProcessKey(Key{Kind: KeyBackspace}); err != nil {
		t.Fatalf("backspace: %v", err)
	}
	line, _ = ed.buffer.Line(0)
	if string(line) != "" {
		t.Fatalf("expected buffer to be empty, got %q", string(line))
	}
}

func TestEditorCursorMovement(t *testing.T) {
	ed := newTestEditor(t)
	ed.buffer = NewBuffer([]string{"one", "two"})
	ed.ProcessKey(Key{Kind: KeyArrowDown})
	if ed.cursor.Row != 1 {
		t.Fatalf("expected cursor row 1, got %d", ed.cursor.Row)
	}
	ed.ProcessKey(Key{Kind: KeyArrowRight})
	if ed.cursor.Col != 1 {
		t.Fatalf("expected cursor col 1, got %d", ed.cursor.Col)
	}
	ed.ProcessKey(Key{Kind: KeyHome})
	if ed.cursor.Col != 0 {
		t.Fatalf("expected cursor col 0, got %d", ed.cursor.Col)
	}
	ed.ProcessKey(Key{Kind: KeyEnd})
	if ed.cursor.Col != len("two") {
		t.Fatalf("expected cursor end, got %d", ed.cursor.Col)
	}
}

func TestEditorSaveAndOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(path, []byte("line1\nline2"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	ed := newTestEditor(t)
	if err := ed.OpenFile(path); err != nil {
		t.Fatalf("open file: %v", err)
	}
	if ed.buffer.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", ed.buffer.LineCount())
	}
	ed.cursor.Row = 1
	ed.cursor.Col = 5
	if err := ed.ProcessKey(RuneKey('!')); err != nil {
		t.Fatalf("insert rune: %v", err)
	}
	if err := ed.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	updated, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(updated) != "line1\nline2!" {
		t.Fatalf("unexpected content: %q", string(updated))
	}
}

func TestEditorRefreshScreen(t *testing.T) {
	output := &bytes.Buffer{}
	ed, err := New(fakeTerminal{rows: 6, cols: 20}, bytes.NewBuffer(nil), output)
	if err != nil {
		t.Fatalf("new editor: %v", err)
	}
	ed.buffer = NewBuffer([]string{"hello"})
	if err := ed.RefreshScreen(); err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if output.Len() == 0 {
		t.Fatalf("expected output to be written")
	}
}
