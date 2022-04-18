package lexer

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	Lexeme       string
	lexemeStack  []string
	tokenStack   []Token
	peek         byte
	Line         uint32
	reader       *bufio.Reader
	read_pointer int
	key_words    map[string]Token
}

func NewLexer(source string) Lexer {
	str := strings.NewReader(source)
	source_reader := bufio.NewReaderSize(str, len(source))
	lexer := Lexer{
		Line:      uint32(1),
		reader:    source_reader,
		key_words: make(map[string]Token),
	}

	lexer.reserve()

	return lexer
}

func (l *Lexer) ReverseScan() {
	/*
		back_len := len(l.Lexeme)
		只能un read 一个字节
		for i := 0; i < back_len; i++ {
			l.reader.UnreadByte()
		}
	*/
	if l.read_pointer > 0 {
		l.read_pointer = l.read_pointer - 1
	}

}

func (l *Lexer) reserve() {
	key_words := GetKeyWords()
	for _, key_word := range key_words {
		l.key_words[key_word.ToString()] = key_word.Tag
	}
}

func (l *Lexer) Readch() error {
	char, err := l.reader.ReadByte() //提前读取下一个字符
	l.peek = char
	return err
}

func (l *Lexer) ReadCharacter(c byte) (bool, error) {

	chars, err := l.reader.Peek(1)
	if err != nil {
		return false, err
	}

	peekChar := chars[0]
	if peekChar != c {
		return false, nil
	}

	l.Readch() //越过当前peek的字符
	return true, nil
}

func (l *Lexer) UnRead() error {
	return l.reader.UnreadByte()
}

func (l *Lexer) Scan() (Token, error) {

	if l.read_pointer < len(l.lexemeStack) {
		l.Lexeme = l.lexemeStack[l.read_pointer]
		token := l.tokenStack[l.read_pointer]
		l.read_pointer = l.read_pointer + 1
		return token, nil
	} else {
		l.read_pointer = l.read_pointer + 1
	}

	for {
		err := l.Readch()
		if err != nil {
			return NewToken(ERROR), err
		}

		if l.peek == ' ' || l.peek == '\t' {
			continue
		} else if l.peek == '\n' {
			l.Line = l.Line + 1
		} else {
			break
		}
	}

	l.Lexeme = ""

	switch l.peek {
	case ';':
		l.Lexeme = ";"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(SEMICOLON)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '{':
		l.Lexeme = "{"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(LEFT_BRACE)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '}':
		l.Lexeme = "}"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(RIGHT_BRACE)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '+':
		l.Lexeme = "+"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(PLUS)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '-':
		l.Lexeme = "-"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(MINUS)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '(':
		l.Lexeme = "("
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(LEFT_BRACKET)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case ')':
		l.Lexeme = ")"
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(RIGHT_BRACKET)
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	case '&':
		l.Lexeme = "&"
		if ok, err := l.ReadCharacter('&'); ok {
			l.Lexeme = "&&"
			word := NewWordToken("&&", AND)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(AND_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}
	case '|':
		l.Lexeme = "|"
		if ok, err := l.ReadCharacter('|'); ok {
			l.Lexeme = "||"
			word := NewWordToken("||", OR)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(OR_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '=':
		l.Lexeme = "="
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "=="
			word := NewWordToken("==", EQ)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(ASSIGN_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '!':
		l.Lexeme = "!"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "!="
			word := NewWordToken("!=", NE)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(NEGATE_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '<':
		l.Lexeme = "<"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = "<="
			word := NewWordToken("<=", LE)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(LESS_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	case '>':
		l.Lexeme = ">"
		if ok, err := l.ReadCharacter('='); ok {
			l.Lexeme = ">="
			word := NewWordToken(">=", GE)
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, word.Tag)
			return word.Tag, err
		} else {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(GREATER_OPERATOR)
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}

	}

	if unicode.IsNumber(rune(l.peek)) {
		var v int
		var err error
		for {
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				if l.peek != 0 { //l.peek == 0 意味着已经读完所有字符
					l.UnRead() //将字符放回以便下次扫描
				}

				break
			}
			v = 10*v + num
			l.Lexeme += string(l.peek)
			l.Readch()
		}

		if l.peek != '.' {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			token := NewToken(NUM)
			token.lexeme = l.Lexeme
			l.tokenStack = append(l.tokenStack, token)
			return token, err
		}
		l.Lexeme += string(l.peek)
		l.Readch() //越过 "."

		x := float64(v)
		d := float64(10)
		for {
			l.Readch()
			num, err := strconv.Atoi(string(l.peek))
			if err != nil {
				if l.peek != 0 { //l.peek == 0 意味着已经读完所有字符
					l.UnRead() //将字符放回以便下次扫描
				}

				break
			}

			x = x + float64(num)/d
			d = d * 10
			l.Lexeme += string(l.peek)
		}
		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token := NewToken(REAL)
		token.lexeme = l.Lexeme
		l.tokenStack = append(l.tokenStack, token)
		return token, err
	}

	if unicode.IsLetter(rune(l.peek)) {
		var buffer []byte
		for {
			buffer = append(buffer, l.peek)
			l.Lexeme += string(l.peek)

			l.Readch()
			if !unicode.IsLetter(rune(l.peek)) {
				if l.peek != 0 { //l.peek == 0 意味着已经读完所有字符
					l.UnRead() //将字符放回以便下次扫描
				}
				break
			}
		}

		s := string(buffer)
		token, ok := l.key_words[s]
		if ok {
			l.lexemeStack = append(l.lexemeStack, l.Lexeme)
			l.tokenStack = append(l.tokenStack, token)
			return token, nil
		}

		l.lexemeStack = append(l.lexemeStack, l.Lexeme)
		token = NewToken(ID)
		token.lexeme = l.Lexeme
		l.tokenStack = append(l.tokenStack, token)
		return token, nil
	}

	return NewToken(EOF), nil
}
