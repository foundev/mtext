package editor

import (
	"bytes"
	"testing"
)

func TestKeyReader(t *testing.T) {
	input := []byte{ctrlQ, ctrlS, 127, '\r', esc, '[', 'A', esc, '[', '3', '~', 'x'}
	kr := NewKeyReader(bytes.NewReader(input))

	cases := []struct {
		expected KeyKind
	}{
		{KeyCtrlQ},
		{KeyCtrlS},
		{KeyBackspace},
		{KeyEnter},
		{KeyArrowUp},
		{KeyDelete},
		{KeyRune},
	}

	for i, c := range cases {
		key, err := kr.ReadKey()
		if err != nil {
			t.Fatalf("read key %d: %v", i, err)
		}
		if key.Kind != c.expected {
			t.Fatalf("expected %v, got %v", c.expected, key.Kind)
		}
	}
}
