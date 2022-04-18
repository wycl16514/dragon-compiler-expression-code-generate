在前面章节中我们给出了语法解析树对应节点的设计，这些节点能够针对其内容完成中间代码的输出，这一节我们继续完善必要节点的设计，然后手动构造语法树，并驱动语法树实现中间代码生成。

首先我们增加一个赋值节点，也就是Set节点的实现，它对应类似赋值语句"c=a;",在inter内添加一个set.go文件，然后添加代码如下：
```
package inter

/*
Set 节点对应 c = a+b,因此它包含两部分，分别是左边的ID节点和右边的Op节点
*/
type Set struct {
	id   ExprInterface
	expr ExprInterface
}

func checkType(p1 *Type, p2 *Type) *Type {
	//c = a + b , c的类型会转换为右边a+b的类型
	if Numberic(p1) && Numberic(p2) {
		return p2
	} else if p1.Lexeme == "bool" && p2.Lexeme == "bool" {
		return p2
	}

	return nil
}

func NewSet(id ExprInterface, expr ExprInterface) (*Set, error) {
	if checkType(id.Type(), expr.Type()) == nil {
		return nil, id.Errors("type error")
	} else {
		return &Set{
			id:   id,
			expr: expr,
		}, nil
	}
}

func (s *Set) Errors(str string) error {
	return s.id.Errors(str)
}

func (s *Set) NewLabel() uint32 {
	return s.id.NewLabel()
}

func (s *Set) EmitLable(i uint32) {
	s.id.EmitLabel(i)
}

func (s *Set) Emit(code string) {
	s.id.Emit(code)
}

func (s *Set) Gen() ExprInterface {
	s.expr = s.expr.Gen()
	s.Emit(s.id.ToString() + " = " + s.expr.ToString())
	return s.id
}

func (s *Set) Reduce() ExprInterface {
	return s.id.Reduce()
}

func (s *Set) Type() *Type {
	return s.id.Type()
}

func (s *Set) ToString() string {
	return s.id.ToString()
}

```
有了赋值节点后，我们就可以针对赋值语句例如"a=b+c"来生成中间代码，此外我们还需要再增加一个节点也就是常量节点，当编译器读取到类似“3;","5;"等常量时就会构造对应节点，在inter下创建文件constant.go，添加代码如下：
```
package inter

import (
	"lexer"
)

type Constant struct {
	expr *Expr 
}

func NewConstant(line uint32, token *lexer.Token, expr_type *Type) *Constant{
	constant := &Constant{
		expr: NewExpr(line, token, expr_type),
	}

	return constant 
}

//定义两个常量 true和false
func GetConstantTrue() *Constant {
	tok := lexer.NewToken(lexer.TRUE)
	true_type := NewType("bool", lexer.TRUE, 1)
	return NewConstant(0, &tok, true_type)
}

func GetConstantFalse() *Constant {
	tok := lexer.NewToken(lexer.FALSE)
	false_type := NewType("bool", lexer.FALSE, 1)
	return NewConstant(0, &tok, false_type)
}

func (c *Constant) Errors(s string) error {
	return c.expr.Errors(s)
}

func (c *Constant) NewLabel() uint32 {
	return c.expr.NewLabel()
}

func (c *Constant) EmitLabel(l uint32) {
	c.expr.EmitLabel(l)
}

func (c *Constant) Emit(code string) {
	c.expr.Emit(code)
}

func (c *Constant) Gen() ExprInterface {
	return c
}

func (c *Constant) Reduce() ExprInterface {
	return c
}

func (c *Constant) Type() *Type {
	return c.expr.Type()
}

func (c *Constant) ToString() string {
	return c.expr.ToString()
}
```
完成上面节点实现后，我们在main.go中手动构造一个语法解析树，其代码如下：
```
package main

import (
	"lexer"

	"inter"
)

func main() {

	expr_type := inter.NewType("int", lexer.BASIC, 4)
	id_a := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "a"), expr_type)
	id_b := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "b"), expr_type)
	//a+b
	arith1, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), id_a, id_b)

	id_c := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "c"), expr_type)
	id_d := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "d"), expr_type)
	arith2, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), id_c, id_d)

	arith3, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.MINUS, "-"), arith1, arith2)

	//arith3.Reduce()

	id_e := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "e"), expr_type)
	//e = (a+b) - (b+c) -> c = t1 - t2
	set, _ := inter.NewSet(id_e, arith3)
	set.Gen()
}
```
上面代码构造的语法树结构如下：
![请添加图片描述](https://img-blog.csdnimg.cn/405696dad6ca4551b429a48c3a22a661.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBAdHlsZXJfZG93bmxvYWQ=,size_18,color_FFFFFF,t_70,g_se,x_16)
首先代码段：
```
expr_type := inter.NewType("int", lexer.BASIC, 4)
id_a := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "a"), expr_type)
id_b := inter.NewID(1, lexer.NewTokenWithString(lexer.ID, "b"), expr_type)
```
手动构造了两个ID节点，分别对应变量a,b，然后代码：
```
arith1, _ := inter.NewArith(1, lexer.NewTokenWithString(lexer.PLUS, "+"), id_a, id_b)
```
将节点a,b和符号"+"合起来形成一个算术表达式，这个算术表达式执行Reduce操作时，会将a和b两个变量相加，然后OP节点会分配一个临时寄存器变量t1，将他们相加的结果存放到t1,同理变量c+d也会通过Reduce操作，将他们相加的结果存放到临时寄存器变量t2,最后t1和t2通过”-“结合成一个算术表达式，与ID节点e一起构成一个Set节点，其中ID节点e对应Set节点的id字段，t1-t2对应Set节点的expr字段，于是在Set节点调用Gen生成代码是就会形成e = t1 - t2的结果。

上面代码运行后所得结果如下：
![请添加图片描述](https://img-blog.csdnimg.cn/575d54bce38141369df62065a2c0bd81.png?x-oss-process=image/watermark,type_d3F5LXplbmhlaQ,shadow_50,text_Q1NETiBAdHlsZXJfZG93bmxvYWQ=,size_20,color_FFFFFF,t_70,g_se,x_16)

可以看到运行结果跟我们的推导是一样的，要想更好的理解代码逻辑，最好还是通过观看调试演示视频，请在b站搜索：Coding迪斯尼，[代码下载地址](https://github.com/wycl16514/dragon-compiler-expression-code-generate.git)：https://github.com/wycl16514/dragon-compiler-expression-code-generate.git，网盘下载为：链接: https://pan.baidu.com/s/1g8VImSml68jEZuuGXpPKog 提取码: 2185，[更多干货](http://m.study.163.com/provider/7600199/index.htm?share=2&shareId=7600199)
