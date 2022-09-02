package fields

import (
	"bufio"
	"strings"
)

type Field struct {
	Parent   *Field
	Children []*Field

	Name string
}

func (f *Field) Has(field string) (bool, *Field) {
	for _, child := range f.Children {
		if child.Name == field {
			return true, child
		}
	}

	return false, nil
}

func (f *Field) ToString() string {
	parts := []string{}

	for _, child := range f.Children {
		currentPart := child.Name
		if len(child.Children) > 0 {
			currentPart += "("
			currentPart += child.ToString()
			currentPart += ")"
		}

		parts = append(parts, currentPart)
	}

	return strings.Join(parts, ",")
}

func ParseFields(fields string) *Field {
	root := &Field{}

	reader := bufio.NewReader(strings.NewReader(fields))

	currentToken := ""
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			// EOF
			node := &Field{Parent: root, Name: currentToken}
			root.Children = append(root.Children, node)
			break
		} else if r == ',' {
			node := &Field{Parent: root, Name: currentToken}
			root.Children = append(root.Children, node)

			currentToken = ""
		} else if r == '(' {
			node := &Field{Parent: root, Name: currentToken}
			root.Children = append(root.Children, node)

			root = node
			currentToken = ""
		} else if r == ')' {
			node := &Field{Parent: root, Name: currentToken}
			root.Children = append(root.Children, node)

			root = root.Parent
			currentToken = ""
		} else {
			currentToken += string(r)
		}
	}

	//log.Printf("%+v", root)
	return root
}
