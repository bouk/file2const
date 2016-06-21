package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type inputFile struct {
	contents     string
	constantName string
}

var (
	packageName string
	outputFile  string
	inputFiles  []inputFile
)

func printHelp() {
	fmt.Fprintf(os.Stderr, "Usage: %s <flags> [<inputFile:constantName>] <outputFile.go>\n", os.Args[0])
	flag.PrintDefaults()
}

func parseFlags() error {
	flag.StringVar(&packageName, "package", "", "name of package to give to file. Default to directory name")
	flag.Parse()
	if flag.NArg() < 2 {
		return fmt.Errorf("need at least 2 arguments")
	}

	outputFile = flag.Args()[flag.NArg()-1]
	outputFile, err := filepath.Abs(outputFile)
	if err != nil {
		return err
	}
	if packageName == "" {
		packageName = filepath.Base(filepath.Dir(outputFile))
	}
	inputFiles = make([]inputFile, flag.NArg()-1)
	for i, input := range flag.Args()[:flag.NArg()-1] {
		s := strings.SplitN(input, ":", 2)
		if len(s) != 2 {
			return fmt.Errorf("wrongly formatted input file %#v", input)
		}
		contents, err := ioutil.ReadFile(s[0])
		if err != nil {
			return err
		}
		inputFiles[i].contents = string(contents)
		inputFiles[i].constantName = s[1]
	}
	return nil
}

func ValueToLiteral(value string) string {
	if strings.Contains(value, "`") || strings.IndexFunc(value, func(r rune) bool { return !unicode.IsGraphic(r) && r != '\n' }) != -1 {
		return strconv.Quote(value)
	} else {
		return "`" + value + "`"
	}
}

func main() {
	err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr)
		printHelp()
		return
	}

	fs := token.NewFileSet()

	values := make([]ast.Spec, len(inputFiles))
	for i, input := range inputFiles {
		values[i] = &ast.ValueSpec{
			Names:  []*ast.Ident{&ast.Ident{Name: input.constantName}},
			Values: []ast.Expr{&ast.BasicLit{Kind: token.STRING, Value: ValueToLiteral(input.contents)}},
		}
	}
	f := &ast.File{
		Name: &ast.Ident{Name: packageName},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:    token.CONST,
				Lparen: token.Pos(1),
				Rparen: token.Pos(1),
				Specs:  values,
			},
		},
	}

	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	defer file.Close()

	format.Node(file, fs, f)
}
