package evaluator

import (
	"interrupter/ast"
	"interrupter/object"
	"interrupter/xlog"
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Boolean:
		return object.TrueOrFase(n.Value)
	case *ast.IntegerLiteral:
		return &object.IntegerObject{Value: n.Value}
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpr(n.Operator, right)
	case *ast.InfixExpression:
		left, right := Eval(n.Left), Eval(n.Right)
		return evalInfixExpr(n.Operator, left, right)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.Program:
		return evalStatements(n.Statements)
	}
	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	xlog.Debugf("eval statements: %#v\n", stmts)
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefixExpr(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalPrefixBangExpr(right)
	case "-":
		return evalPrefixSubExpr(right)
	}
	return nil
}

func evalInfixExpr(op string, left, right object.Object) object.Object {
	switch op {
	case "+":
		return evalInfixPlusExpr(left, right)
	case "-":
		return evalInfixSubExpr(left, right)
	case "*":
		return evalInfixMultiExpr(left, right)
	case "/":
		return evalInfixDivExpr(left, right)
	case ">":
		return evalInfixGTExpr(left, right)
	case "<":
		return evalInfixLTExpr(left, right)
	case "==":
		return evalInfixEQTExpr(left, right)
	case "!=":
		return evalInfixNOTEQTExpr(left, right)
	default:
		return object.NULL
	}
}

func evalInfixPlusExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.IntegerObject{Value: l.Value + r.Value}
	}
	return object.NULL
}

func evalInfixSubExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.IntegerObject{Value: l.Value - r.Value}
	}
	return object.NULL
}

func evalInfixMultiExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.IntegerObject{Value: l.Value * r.Value}
	}
	return object.NULL
}

func evalInfixDivExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.IntegerObject{Value: l.Value / r.Value}
	}
	return object.NULL
}

func evalInfixGTExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.BooleanObject{Value: l.Value > r.Value}
	}
	return object.NULL
}

func evalInfixLTExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.BooleanObject{Value: l.Value < r.Value}
	}
	return object.NULL
}
func evalInfixEQTExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.BooleanObject{Value: l.Value == r.Value}
	}
	return object.NULL
}
func evalInfixNOTEQTExpr(left, right object.Object) object.Object {
	l, lOk := left.(*object.IntegerObject)
	r, rOk := right.(*object.IntegerObject)
	if lOk && rOk {
		return &object.BooleanObject{Value: l.Value != r.Value}
	}
	return object.NULL
}

func evalPrefixBangExpr(obj object.Object) object.Object {
	switch obj {
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
	}
}

func evalPrefixSubExpr(obj object.Object) object.Object {
	switch o := obj.(type) {
	case *object.IntegerObject:
		return &object.IntegerObject{Value: -o.Value}
	default:
		return object.NULL
	}
}
