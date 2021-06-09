package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"text/template"
)

type TypeField struct {
	GoName  string
	SqlName string
	Type    string
}

const storeTemplate = `
type {{.Typename}}Store struct {
	db *sql.DB
}

func (s *{{.Typename}}Store) Get(ids []int) ([]*{{.Typename}}, error) {
	query := "select {{.SqlFields}} from {{.TableName}} where id in (" + idsStr(ids) + ")"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*{{.Typename}}, 0)
	for rows.Next() {
		obj := {{.Typename}}{}
		err := rows.Scan({{.ScanFields}})
		if err != nil {
			return nil, err
		}
		result = append(result, &obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
`

type TemplateVars struct {
	Typename   string
	SqlFields  string
	TableName  string
	ScanFields string
}

func main() {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "./", nil, parser.AllErrors|parser.ParseComments)
	log.Printf("%v", pkgs)
	log.Printf("%v", err)

	mainPkg, ok := pkgs["main"]
	if !ok {
		log.Fatalf("main package not found")
	}
	storeFile, ok := mainPkg.Files["store.go"]
	if !ok {
		log.Fatalf("store.go file not found")
	}

	resultFile := "package main\n\nimport \"database/sql\"\n"
	stores := make([]string, 0)

	for _, decl := range storeFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			fields := make([]TypeField, 0)
			log.Printf("Type: %s", typeSpec.Name.Name)
			for _, field := range structType.Fields.List {
				fieldName := field.Names[0].Name
				if field.Tag == nil || field.Tag.Value == "" {
					continue
				}

				fieldType := field.Type.(*ast.Ident).Name
				fieldTag := field.Tag.Value
				tagRegexp := regexp.MustCompile(`db:"[a-z0-9_]+"`)
				fieldDbName := tagRegexp.FindString(fieldTag)
				fieldDbName = fieldDbName[4 : len(fieldDbName)-1]
				log.Printf("  -> %s:%s %s", fieldType, fieldName, fieldDbName)

				fields = append(fields, TypeField{
					GoName:  fieldName,
					SqlName: fieldDbName,
					Type:    fieldType,
				})
			}

			if len(fields) == 0 {
				continue
			}

			sqlFields := make([]string, len(fields))
			scanFields := make([]string, len(fields))
			for i, field := range fields {
				if field.Type == "int" {
					sqlFields[i] = "coalesce(" + field.SqlName + ", 0)"
				} else if field.Type == "string" {
					sqlFields[i] = "coalesce(" + field.SqlName + ", '')"
				} else {
					log.Fatalf("Unknown type: %s", field.Type)
				}

				scanFields[i] = "&obj." + field.GoName
			}

			storeTemplateParser, err := template.New("store").Parse(storeTemplate)
			if err != nil {
				log.Fatal(err)
			}

			w := bytes.NewBuffer(nil)
			err = storeTemplateParser.Execute(w, TemplateVars{
				Typename:   typeSpec.Name.Name,
				SqlFields:  strings.Join(sqlFields, ", "),
				TableName:  strings.ToLower(typeSpec.Name.Name),
				ScanFields: strings.Join(scanFields, ", "),
			})
			if err != nil {
				log.Fatal(err)
			}

			resultFile += w.String()
			stores = append(stores, typeSpec.Name.Name)
		}
	}

	resultFile += "type AllStores struct {\n"
	for _, store := range stores {
		resultFile += "    " + store + " *" + store + "Store\n"
	}
	resultFile += "}\nfunc NewAllStores(db *sql.DB) *AllStores {\n"
	resultFile += "    return &AllStores {\n"
	for _, store := range stores {
		resultFile += "    " + store + ": &" + store + "Store{db: db},\n"
	}
	resultFile += "    }\n}\n"

	err = ioutil.WriteFile("store_gen.go", []byte(resultFile), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
