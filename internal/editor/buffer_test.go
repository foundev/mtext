package editor

import "testing"

func TestBufferInsertDeleteSplitMerge(t *testing.T) {
	buf := NewBuffer([]string{"hello"})
	if err := buf.InsertRune(0, 5, '!'); err != nil {
		t.Fatalf("insert rune: %v", err)
	}
	line, _ := buf.Line(0)
	if string(line) != "hello!" {
		t.Fatalf("expected inserted line, got %q", string(line))
	}
	if err := buf.DeleteRune(0, 5); err != nil {
		t.Fatalf("delete rune: %v", err)
	}
	line, _ = buf.Line(0)
	if string(line) != "hello" {
		t.Fatalf("expected delete result, got %q", string(line))
	}
	if err := buf.SplitLine(0, 2); err != nil {
		t.Fatalf("split line: %v", err)
	}
	if buf.LineCount() != 2 {
		t.Fatalf("expected 2 lines, got %d", buf.LineCount())
	}
	line, _ = buf.Line(0)
	if string(line) != "he" {
		t.Fatalf("expected split left, got %q", string(line))
	}
	line, _ = buf.Line(1)
	if string(line) != "llo" {
		t.Fatalf("expected split right, got %q", string(line))
	}
	if err := buf.MergeLine(0); err != nil {
		t.Fatalf("merge line: %v", err)
	}
	if buf.LineCount() != 1 {
		t.Fatalf("expected merged line count, got %d", buf.LineCount())
	}
	line, _ = buf.Line(0)
	if string(line) != "hello" {
		t.Fatalf("expected merge result, got %q", string(line))
	}
}
