// ************************************************************************************************
// Assembler

type Assembler struct {
    input string
    runes []rune
    index int
}

// TODO: Flesh out the lexer. Should it be a separate object?

func (asm *Assembler) getNextChar() rune

func (asm *Assembler) getNextToken() (*Token) {
    if asm.index >= len(asm.runes) {
        return Token{"EOF", TokEOF}
    }
    ch := asm.runes[asm.index++]
    if unicode.IsDigit(ch) {
        // parse a number
        tokstr := []runes{ch}
        
    i := asm.index
    for _, v := range asm.input[i:] {
        if !unicode.IsSpace(v) {
            break
        }
        i++
    }
    if i >= len(asm.input) {
        return "", true
    }
    j := i
    for _, v := range(asm.input[j:]) {
        if unicode.IsSpace(v) {
            break
        }
        j++
    }
    asm.index = j
    return asm.input[i:j], false
}

// TODO: Do this.
func Assemble(input string) ([]uint16, error) {
    asm := Assembler{input, []rune(input), 0}
    var result []uint16
    for tok, eof := asm.getNextToken(); !eof; tok, eof = asm.getNextToken() {
        // switch tok {
        // case "SETR":
        //     // TODO: Fill this in.
    }
    return result, nil
}
