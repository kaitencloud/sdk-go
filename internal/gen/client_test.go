// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package gen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"testing"
)

// Raw exposes the generated client through the public SDK. Check declarations,
// not just its method set: request builders, response parsers and body aliases
// must not provide a second route to direct flag evaluation. Protocol models
// and the feature-flag management/webhook operations remain legitimate.
func TestNoDirectFlagEvaluationExports(t *testing.T) {
	t.Parallel()

	file, err := parser.ParseFile(token.NewFileSet(), "client.gen.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	ast.Inspect(file, func(node ast.Node) bool {
		if literal, ok := node.(*ast.BasicLit); ok && literal.Kind == token.STRING {
			value, err := strconv.Unquote(literal.Value)
			if err != nil {
				t.Error(err)
			} else if strings.Contains(value, "/ofrep/v1/evaluate/flags") {
				t.Error("forbidden direct evaluation request path")
			}
		}
		ident, ok := node.(*ast.Ident)
		if ok && ident.IsExported() && (strings.Contains(ident.Name, "EvaluateFlag") || strings.Contains(ident.Name, "OfrepEvaluate")) {
			t.Errorf("forbidden direct evaluation symbol: %s", ident.Name)
		}
		return true
	})
}
