package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
	"log"
	"os"
	"strings"
)

func Serialize() {
	r := UsersFollowReq{
		TargetId: "",
		Action:   SubscribeAction_FOLLOW,
	}
	resp, err := json.Marshal(r)
	log.Printf("%s %s", resp, err)
}

func DoParse() error {
	protoFile, err := os.ReadFile("../schema/api.proto")
	if err != nil {
		return fmt.Errorf("cannot open .proto file: %w", err)
	}

	file, err := protoparser.Parse(bytes.NewBuffer(protoFile))
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	result := "package api\n\n"
	result += "import \"context\"\n\n"

	resultEnums, enumsMap := generateEnums(file.ProtoBody)
	result += resultEnums

	for _, node := range file.ProtoBody {
		switch node := node.(type) {
		case *parser.Service:
			result += generateService(node)
		case *parser.Message:
			result += generateMessage(node, enumsMap)
		}
	}

	log.Printf("%+v", file)
	log.Printf("%s", result)

	_ = os.WriteFile("api.go", []byte(result), 0755)

	return nil
}

func generateEnums(nodes []parser.Visitee) (string, map[string]bool) {
	result := ""
	enumsMap := map[string]bool{}

	for _, node := range nodes {
		if node, ok := node.(*parser.Enum); ok {
			result += "type " + node.EnumName + " int\n\n"
			result += "const (\n"

			for _, enumNodeItem := range node.EnumBody {
				enumField, ok := enumNodeItem.(*parser.EnumField)
				if ok {
					result += "\t" + node.EnumName + "_" + enumField.Ident + " " + node.EnumName + " = " + enumField.Number + "\n"
				}
			}

			result += ")\n\n"
			enumsMap[node.EnumName] = true
		}
	}

	return result, enumsMap
}

type message struct {
	fields []field
}

type field struct {
	Name   string
	Type   string
	Number int
}

func generateMessage(node *parser.Message, enumsMap map[string]bool) string {
	log.Printf("Found message: %s", node.MessageName)

	result := ""
	result += "type " + node.MessageName + " struct {\n"

	for _, node := range node.MessageBody {
		fieldNode, ok := node.(*parser.Field)
		if ok {
			log.Printf("%s %s",
				fieldNode.FieldName,
				fieldNode.Type,
			)

			goType := ""
			if fieldNode.Type == "int32" {
				// Int
				goType = "int"
			} else if fieldNode.Type == "string" || fieldNode.Type == "bool" {
				// Scalars
				goType = fieldNode.Type
			} else if enumsMap[fieldNode.Type] {
				// Enum, this is int scalar
				goType = fieldNode.Type
			} else {
				// Pointer
				goType = "*" + fieldNode.Type
			}

			fieldType := ""
			if fieldNode.IsRepeated {
				fieldType = "[]"
			}
			fieldType += goType

			result += "\t" + strings.Title(fieldNode.FieldName)
			result += " " + fieldType
			result += " `json:\"" + fieldNode.FieldName + ",omitempty\"`"
			result += "\n"
		}
	}

	result += "}\n\n"

	return result
}

func generateService(node *parser.Service) string {
	result := ""
	result += "type " + node.ServiceName + " interface {\n"

	for _, node := range node.ServiceBody {
		rpcNode, ok := node.(*parser.RPC)
		if ok {
			result += "\t" + rpcNode.RPCName + "(context.Context, *" + rpcNode.RPCRequest.MessageType + ") (*" + rpcNode.RPCResponse.MessageType + ", error)\n"
		}
	}

	result += "}\n\n"

	return result
}
