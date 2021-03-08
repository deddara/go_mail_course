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
	"strings"
)

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
	fmt.Println(mainStruct)
}
