package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

var (
	TRUE  = &BooleanObject{Value: true}
	FALSE = &BooleanObject{Value: false}
	NULL  = &NullObject{}
)

func TrueOrFase(isTrue bool) *BooleanObject {
	if isTrue {
		return TRUE
	}
	return FALSE
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type IntegerObject struct {
	Value int64
}

func (i *IntegerObject) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *IntegerObject) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type BooleanObject struct {
	Value bool
}

func (b *BooleanObject) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *BooleanObject) Inspect() string {
	return fmt.Sprintf("%v", b.Value)
}

type NullObject struct{}

func (n *NullObject) Type() ObjectType {
	return NULL_OBJ
}

func (n *NullObject) Inspect() string {
	return "null"
}
