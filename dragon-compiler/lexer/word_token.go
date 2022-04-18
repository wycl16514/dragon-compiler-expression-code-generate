package lexer

type Word struct {
	lexeme string
	Tag    Token
}

func NewWordToken(s string, tag Tag) Word {
	return Word{
		lexeme: s,
		Tag:    NewToken(tag),
	}
}

func (w *Word) ToString() string {
	return w.lexeme
}

func GetKeyWords() []Word {
	key_words := []Word{}
	key_words = append(key_words, NewWordToken("&&", AND))
	key_words = append(key_words, NewWordToken("||", OR))
	key_words = append(key_words, NewWordToken("==", EQ))
	key_words = append(key_words, NewWordToken("!=", NE))
	key_words = append(key_words, NewWordToken("<=", LE))
	key_words = append(key_words, NewWordToken(">=", GE))
	key_words = append(key_words, NewWordToken("minus", MINUS))
	key_words = append(key_words, NewWordToken("true", TRUE))
	key_words = append(key_words, NewWordToken("false", FALSE))
	key_words = append(key_words, NewWordToken("if", IF))
	key_words = append(key_words, NewWordToken("else", ELSE))
	//添加类型定义
	key_words = append(key_words, NewWordToken("int", BASIC))
	key_words = append(key_words, NewWordToken("float", BASIC))
	key_words = append(key_words, NewWordToken("bool", BASIC))
	key_words = append(key_words, NewWordToken("char", BASIC))

	return key_words
}
