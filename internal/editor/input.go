package editor

import (
	"bufio"
	"errors"
	"io"
)

const (
	ctrlQ = 17
	ctrlS = 19
	esc   = 27
)

type KeyReader struct {
	reader *bufio.Reader
}

func NewKeyReader(input io.Reader) *KeyReader {
	return &KeyReader{reader: bufio.NewReader(input)}
}

func (kr *KeyReader) ReadKey() (Key, error) {
	b, err := kr.reader.ReadByte()
	if err != nil {
		return Key{}, err
	}

	switch b {
	case ctrlQ:
		return Key{Kind: KeyCtrlQ}, nil
	case ctrlS:
		return Key{Kind: KeyCtrlS}, nil
	case 127:
		return Key{Kind: KeyBackspace}, nil
	case '\r':
		return Key{Kind: KeyEnter}, nil
	case esc:
		return kr.readEscapeSequence()
	default:
		if b < 32 {
			return Key{Kind: KeyUnknown}, nil
		}
		return RuneKey(rune(b)), nil
	}
}

func (kr *KeyReader) readEscapeSequence() (Key, error) {
	first, err := kr.reader.ReadByte()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return Key{Kind: KeyUnknown}, nil
		}
		return Key{}, err
	}
	if first != '[' && first != 'O' {
		return Key{Kind: KeyUnknown}, nil
	}

	second, err := kr.reader.ReadByte()
	if err != nil {
		return Key{}, err
	}

	if first == 'O' {
		switch second {
		case 'H':
			return Key{Kind: KeyHome}, nil
		case 'F':
			return Key{Kind: KeyEnd}, nil
		}
		return Key{Kind: KeyUnknown}, nil
	}

	if second >= '0' && second <= '9' {
		third, err := kr.reader.ReadByte()
		if err != nil {
			return Key{}, err
		}
		if third != '~' {
			return Key{Kind: KeyUnknown}, nil
		}
		switch second {
		case '1', '7':
			return Key{Kind: KeyHome}, nil
		case '4', '8':
			return Key{Kind: KeyEnd}, nil
		case '3':
			return Key{Kind: KeyDelete}, nil
		case '5':
			return Key{Kind: KeyPageUp}, nil
		case '6':
			return Key{Kind: KeyPageDown}, nil
		}
		return Key{Kind: KeyUnknown}, nil
	}

	switch second {
	case 'A':
		return Key{Kind: KeyArrowUp}, nil
	case 'B':
		return Key{Kind: KeyArrowDown}, nil
	case 'C':
		return Key{Kind: KeyArrowRight}, nil
	case 'D':
		return Key{Kind: KeyArrowLeft}, nil
	case 'H':
		return Key{Kind: KeyHome}, nil
	case 'F':
		return Key{Kind: KeyEnd}, nil
	}

	return Key{Kind: KeyUnknown}, nil
}
