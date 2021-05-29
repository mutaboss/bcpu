// ************************************************************************************************
// Lexer

type TokenType uint16

const (
	TokHalt   TokenType = iota
	TokNoop
	TokJmp
	TokJeq
	TokJgt
	TokJlt
	TokSetReg
	TokLoad
	TokStore
	TokAddReg
	TokSubReg
	TokMulReg
	TokDivReg
	TokCmp
    TokAnd
    TokOr
    TokXor
    TokShl
    TokShr
    TokNot
    TokInteger
    TokEOF
)

type Token struct {
    literal string
    token TokenType
}

type Lexer struct {
    input []rune
    position int
    readPosition int
    ch    rune
    line int
}

func (l *Lexer) NewLexer(input string) *Lexer {
    l := &Lexer{input: input}
    l.readChar()
    return l
}

func (l *Lexer) readChar() {
    if l.ch == '\n' {
        l.line += 1
    }
    if l.readPosition >= len(l.input) {
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
        l.position = l.readPosition
        l.readPosition += 1
    }
}

func (l *Lexer) NextToken() *Token {
    var tok Token
    if l.ch == 0 {
        return &Token{"EOF", TokEOF}
    } else if unicode.IsDigit(l.ch) {
        return l.readNumber()
    } else if unicode.IsLetter(l.ch) {
        return l.readIdentifier()
    } else {
        fmt.Printf("Invalid character line %d: %s.\n", l.line, l.ch)
        l.readChar()
    }
    return nil // We should not be here!
}

func (l *Lexer) readNumber() *Token {
    var result []rune
    for unicode.IsDigit(l.ch) {
        result := append(result, l.ch)
        l.readChar()
    }
    return &Token{string(result), TokNumber}
}

func (l *Lexer) readIdentifier() *Token {
    var result []rune
    for unicode.IsLetter(l.ch) {
        result := append(result, l.ch)
        l.readChar()
    }
    if ident, err := l.lookupIdentifier(string(result)); err != nil {
        fmt.
    return &Token{string(result), l.lookupIdentifier(string(result))}
}

func (l *Lexer) lookupIdentifier(Identifier) TokenType {
}

