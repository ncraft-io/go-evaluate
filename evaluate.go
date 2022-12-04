// Copyright 2022 ncraft.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evaluate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const GoEvaluate = "//go:evaluate"

type Evaluation struct {
	filename string
	line     int

	fileSet *token.FileSet
}

func New() *Evaluation {
	return &Evaluation{}
}

func (v *Evaluation) Evaluate(filename string) error {
	// os.Open(filename)
	// fs.Stat()

	v.filename = filename
	v.fileSet = token.NewFileSet()
	if ast, err := parser.ParseFile(v.fileSet, filename, nil, parser.ParseComments); err != nil {
		return err
	} else {
		if err = v.evaluate(ast); err != nil {
			return err
		}

		code := bytes.NewBuffer(nil)
		if err = printer.Fprint(code, v.fileSet, ast); err != nil {
			return err
		}

		if err = ioutil.WriteFile(filename, code.Bytes(), fs.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (v *Evaluation) evaluate(file *ast.File) error {
	for _, d := range file.Decls {
		switch x := d.(type) {
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				if decl, ok := spec.(*ast.ValueSpec); ok {
					commentGroup := v.preComment(file, x.Pos())
					if commentGroup != nil && len(commentGroup.List) > 0 {
						comment := commentGroup.List[len(commentGroup.List)-1]
						command := comment.Text
						if isGoEvaluate(command) {
							if len(decl.Names) == 1 && len(decl.Values) == 1 {
								v.line = v.fileSet.Position(comment.Pos()).Line
								command = strings.TrimPrefix(command, GoEvaluate)
								if err := v.updateValue(command, decl); err != nil {
									return err
								}
							} else {
								return v.errorf("only support single value declaration")
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (v *Evaluation) preComment(file *ast.File, pos token.Pos) *ast.CommentGroup {
	for _, comment := range file.Comments {
		ends := comment.End()
		if pos == ends+1 {
			return comment
		}
	}
	return nil
}

func (v *Evaluation) updateValue(command string, decl *ast.ValueSpec) error {
	val, err := v.getValue(command)
	if err != nil {
		return err
	}

	switch expr := decl.Values[0].(type) {
	case *ast.BasicLit:
		switch expr.Kind {
		case token.STRING:
			decl.Values[0] = &ast.BasicLit{
				Value:    strconv.Quote(val),
				Kind:     token.STRING,
				ValuePos: expr.ValuePos,
			}
		default:
			return v.errorf("not support non string type value")
		}
	default:
		return v.errorf("not support non BasicLiteral expression")
	}
	return nil
}

func (v *Evaluation) getValue(command string) (string, error) {
	words, err := v.split(command)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(words[0], words[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", v.errorf("running %q: %s", strings.Join(words, " "), err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (v *Evaluation) split(line string) ([]string, error) {
	var words []string

process:
	for {
		line = strings.TrimLeft(line, " \t")
		if len(line) == 0 {
			break
		}

		if line[0] == '"' {
			quoteClosed := false
			for i := 1; i < len(line); i++ {
				c := line[i]
				if c == '\\' {
					if i+1 == len(line) {
						return nil, v.errorf("bad backslash")
					}
					i++
				} else if c == '"' {
					word, err := strconv.Unquote(line[0 : i+1])
					if err != nil {
						return nil, v.errorf("bad quoted string")
					}
					words = append(words, word)
					line = line[i+1:]

					if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
						return nil, v.errorf("expect space after quoted argument")
					}
					quoteClosed = true
					continue process
				}
			}
			if !quoteClosed {
				return nil, v.errorf("mismatched quoted string")
			}
		}
		i := strings.IndexAny(line, " \t")
		if i < 0 {
			i = len(line)
		}
		words = append(words, line[0:i])
		line = line[i:]
	}

	for i, word := range words {
		words[i] = os.Expand(word, os.Getenv)
	}

	return words, nil
}

func (v *Evaluation) errorf(format string, args ...any) error {
	return fmt.Errorf("%s:%d: %s\n", v.filename, v.line, fmt.Sprintf(format, args...))
}

func isGoEvaluate(buf string) bool {
	return strings.HasPrefix(buf, GoEvaluate+" ") || strings.HasPrefix(buf, GoEvaluate+"\t")
}
