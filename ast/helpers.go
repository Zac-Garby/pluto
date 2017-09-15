package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Zac-Garby/pluto/token"
)

const treeIndent = 2

func in(indent int) string {
	return strings.Repeat(" ", treeIndent*indent)
}

func prefix(indent int, name string) string {
	str := in(indent)

	if name != "" {
		str += name + " ‣ "
	}

	return str
}

// Tree returns a tree representation of a node
func Tree(node Node, indent int, name string) string {
	val := reflect.ValueOf(node)

	typeName := fmt.Sprintf("%T", node)[5:]

	if name != "" {
		name = fmt.Sprintf("%s (%s)", name, typeName)
	} else {
		name = fmt.Sprintf("(%s)", typeName)
	}

	str := prefix(indent, name) + val.Type().Name()

	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Field(i).Interface()
		label := val.Type().Field(i).Name

		if label == "Value" {
			label = ""
		}

		if _, ok := field.(*Null); ok {
			str += "\n" + prefix(indent+1, label) + "<null>"
			continue
		}

		switch n := field.(type) {
		case token.Token:
		case Node:
			str += "\n" + Tree(n, indent+1, label)
		case map[Statement]Statement:
			nodes := make(map[Node]Node)

			for key, value := range n {
				nodes[key] = value
			}

			str += "\n" + makeDictTree(indent+1, nodes, label)
		case map[Expression]Expression:
			nodes := make(map[Node]Node)

			for key, value := range n {
				nodes[key] = value
			}

			str += "\n" + makeDictTree(indent+1, nodes, label)
		case []Statement:
			var nodes []Node

			for _, item := range n {
				nodes = append(nodes, item.(Node))
			}

			str += "\n" + makeListTree(indent+1, nodes, label)
		case []Expression:
			var nodes []Node

			for _, item := range n {
				nodes = append(nodes, item.(Node))
			}

			str += "\n" + makeListTree(indent+1, nodes, label)
		case []Arm:
			str += "\n" + prefix(indent+1, label) + "arms ["

			if len(n) == 0 {
				str += "]"
				break
			}

			for _, arm := range n {
				var nodes []Node

				for _, item := range arm.Exprs {
					nodes = append(nodes, item.(Node))
				}

				str += fmt.Sprintf(
					"\n%sarm (\n%s\n%s\n%s)",
					in(indent+2),
					makeListTree(indent+3, nodes, "expressions"),
					Tree(arm.Body, indent+3, "body"),
					in(indent+2),
				)
			}

			str += "\n" + in(indent+1) + "]"
		case []EmittedItem:
			str += "\n" + prefix(indent+1, label) + "items ["

			if len(n) == 0 {
				str += "]"
				break
			}

			for _, item := range n {
				if item.IsInstruction {
					str += fmt.Sprintf(
						"\n%sinstruction (\n%s%s\n%s%d\n%s)",
						in(indent+2),
						prefix(indent+3, "instruction"),
						item.Instruction,
						prefix(indent+3, "argument"),
						item.Argument,
						in(indent+2),
					)
				} else {
					str += fmt.Sprintf(
						"\n%s",
						Tree(item.Exp, indent+2, "expression"),
					)
				}
			}

			str += "\n" + in(indent+1) + "]"
		default:
			str += "\n" + fmt.Sprintf("%s%s", prefix(indent+1, label), fmt.Sprintf("%v", n))
		}
	}

	return str
}

func makeListTree(indent int, nodes []Node, name string) string {
	str := prefix(indent, name) + "["

	if len(nodes) == 0 {
		return str + "]"
	}

	for _, node := range nodes {
		str += "\n" + Tree(node, indent+1, "")
	}

	return str + "\n" + in(indent) + "]"
}

func makeDictTree(indent int, pairs map[Node]Node, name string) string {
	str := prefix(indent, name) + "["

	if len(pairs) == 0 {
		return str + ":]"
	}

	for key, value := range pairs {
		str += fmt.Sprintf("%s\n%s\n%s\n",
			in(indent),
			Tree(key, indent+1, "key"),
			Tree(value, indent+1, "value"),
		)
	}

	return str + in(indent) + "]"
}
