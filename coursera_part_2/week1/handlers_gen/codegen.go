package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

)

type tmpl struct{
	FieldName string
	HandlerName string
}

// код писать тут


var mainStruct = make(map[string] []methodInfo)

type methodInfo struct {
	structName    string
	structType    ast.Expr

	Url        string
	Method     string
	Auth       bool
}

func checkComment(fDecl *ast.FuncDecl) (bool, methodInfo) {
	comment := fDecl.Doc
	if comment == nil { // if has comment
		log.Printf("SKIP func %#v doesnt have comments\n", fDecl.Name.Name)
		return false, methodInfo{}
	}

	if !strings.HasPrefix(comment.Text(), "apigen:api ") {
		log.Printf("SKIP %#v doent have apigen prefix %s\n", fDecl.Name.Name, comment.Text())
		return false, methodInfo{}
	}
	methodInf := &methodInfo{}
	jsonText := strings.TrimPrefix(comment.Text(), "apigen:api ")
	err := json.Unmarshal([]byte(jsonText), methodInf)
	if err != nil{
		log.Fatal(err.Error())
	}
	return true, *methodInf
}

var (
	ServeHTTPtmpl = template.Must(template.New("ServeHTTPtmpl").Parse(`
// {{.FieldName}}
func (h {{.FieldName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "`))
	methodTmpl = template.Must(template.New("methodTmpl").Parse(`
		h.{{.FieldName}}Handler(w, r)
`))
	handlerTmpl = template.Must(template.New("handlerTmpl").Parse(`
func (h {{.FieldName}}) {{.HandlerName}}Handler(w http.ResponseWriter, r *http.Request) {
`))
)


func genServeHTTP(mainStruct map[string] []methodInfo, out *os.File){
	for key, methods := range mainStruct{
		fmt.Println(key)
		ServeHTTPtmpl.Execute(out, tmpl{key, ""})
		for i, method := range methods {
			out.WriteString(method.Url + "\":")
			methodTmpl.Execute(out, tmpl{method.structName, ""})
			if i + 1 == len(methods){
				out.WriteString(`	default:
		return
	}
`)
			} else{
				out.WriteString("	case \"")
			}
		}
		out.WriteString("}\n")
	}
}

func buildHandler(mainStruct map[string] []methodInfo, out *os.File, node *ast.File){
	//going through all serveHttp structs
	for key, methods := range mainStruct {
		//going through methods of this struct
		for _, method := range methods {
			handlerTmpl.Execute(out, tmpl{key, method.structName})
			//going through decl in file
			for _, f := range node.Decls {
				g, ok := f.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range g.Specs {
					currType, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					//if name of struct is not the same - continue
					if currType.Name.Name != method.structName{
						continue
					}
					currStruct, ok := currType.Type.(*ast.StructType)
					if !ok {
						continue
					}
					// going through all decls in struct
					for _, field := range currStruct.Fields.List {
						if field.Tag == nil {
							continue
						}
						//find apivalidator tag
						tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
						raw_str := tag.Get("apivalidator")
						fmt.Println(raw_str)
					}
				}
			}
		}
	}
}

func main() {
	fSet := token.NewFileSet()
	node, err := parser.ParseFile(fSet, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out, `import "net/http"`)

	for _, f := range node.Decls {
		fDecl, ok := f.(*ast.FuncDecl)
		if !ok {
			log.Printf("SKIP %T is not *ast.FuncDecl\n", f)
			continue
		}
		fRecv := fDecl.Recv
		if fRecv != nil { // if is method
			comment, methodInf := checkComment(fDecl)
			if !comment { //check if has comments and parse comment in json
				continue
			}

			var typeNameBuf bytes.Buffer
			err := printer.Fprint(&typeNameBuf, fSet, fRecv.List[0].Type)
			if err != nil {
				log.Fatalf("failed printing %s", err)
			}
			mainStructName := typeNameBuf.String()
			fmt.Printf("func has struct %s\n", typeNameBuf.String()) //got main struct name type
			//mainStruct[mainStructName] = append(mainStruct[mainStructName], methodInf)

			fIncParams := fDecl.Type.Params.List
			if len(fIncParams) > 1 {
				typeNameBuf.Reset()
				err = printer.Fprint(&typeNameBuf, fSet, fIncParams[1].Type)
				if err != nil {
					log.Fatalf("failed printing %s", err)
				}
				fmt.Printf("incoming type is %s\n", typeNameBuf.String())
				methodInf.structName = typeNameBuf.String()
				methodInf.structType = fIncParams[1].Type
				mainStruct[mainStructName] = append(mainStruct[mainStructName], methodInf)

			}
			fmt.Printf("%#v\n", methodInf)
		}
	}
	genServeHTTP(mainStruct, out)
	buildHandler(mainStruct, out, node)
	fmt.Println(mainStruct)
}
