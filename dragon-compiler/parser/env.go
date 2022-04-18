package simple_parser

type Env struct {
	table map[string]*Symbol
	prev  *Env //这里形成链表
}

func NewEnv(p *Env) *Env {
	return &Env{
		table: make(map[string]*Symbol),
		prev:  p,
	}
}

func (e *Env) Put(s string, sym *Symbol) {
	e.table[s] = sym
}

func (e *Env) Get(s string) *Symbol {
	//查询变量符号时，如果当前符号表没有定义，我们要往上一层作用域做进一步查询
	for env := e; env != nil; env = e.prev {
		found, ok := env.table[s]
		if ok {
			return found
		}
	}

	return nil
}
