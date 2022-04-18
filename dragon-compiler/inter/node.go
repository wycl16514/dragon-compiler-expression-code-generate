package inter

import (
	"errors"
	"fmt"
	"strconv"
)

var labels uint32 //用于实现跳转的标号

type Node struct {
	lex_line uint32
}

func NewNode(line uint32) *Node {
	labels = 0
	return &Node{
		lex_line: line,
	}
}

func (n *Node) Errors(s string) error {
	err_s := "\nnear line " + strconv.FormatUint(uint64(n.lex_line), 10) + s
	return errors.New(err_s)
}

func (n *Node) NewLabel() uint32 {
	labels = labels + 1
	return labels
}

func (n *Node) EmitLabel(i uint32) {
	fmt.Print("\nL" + strconv.FormatUint(uint64(i), 10) + ":\n")
}

func (n *Node) Emit(s string) {
	fmt.Print("\t" + s)
}
