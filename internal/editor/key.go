package editor

type KeyKind int

const (
	KeyUnknown KeyKind = iota
	KeyRune
	KeyCtrlQ
	KeyCtrlS
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyBackspace
	KeyDelete
	KeyEnter
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
)

type Key struct {
	Kind KeyKind
	Rune rune
}

func RuneKey(value rune) Key {
	return Key{Kind: KeyRune, Rune: value}
}
