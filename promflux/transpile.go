package main

import (
	"fmt"
	"time"

	"github.com/influxdata/flux/ast"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
)

type transpiler struct {
	bucket     string
	start      time.Time
	end        time.Time
	resolution time.Duration
}

func buildPipeline(arg ast.Expression, calls ...*ast.CallExpression) *ast.PipeExpression {
	switch len(calls) {
	case 0:
		panic("empty calls list")
	case 1:
		return &ast.PipeExpression{
			Argument: arg,
			Call:     calls[0],
		}
	default:
		return &ast.PipeExpression{
			Argument: buildPipeline(arg, calls[0:len(calls)-1]...),
			Call:     calls[len(calls)-1],
		}
	}
}

func call(fn string, args map[string]ast.Expression) *ast.CallExpression {
	expr := &ast.CallExpression{
		Callee: &ast.Identifier{Name: fn},
	}
	if len(args) > 0 {
		props := make([]*ast.Property, 0, len(args))
		for k, v := range args {
			props = append(props, &ast.Property{
				Key:   &ast.Identifier{Name: k},
				Value: v,
			})
		}

		expr.Arguments = []ast.Expression{
			&ast.ObjectExpression{
				Properties: props,
			},
		}
	}
	return expr
}

// Function to remove extraneous windows at edges of range.
func windowCutoffFn(minStop time.Time, maxStart time.Time) *ast.FunctionExpression {
	return &ast.FunctionExpression{
		Params: []*ast.Property{
			{
				Key: &ast.Identifier{
					Name: "r",
				},
			},
		},
		Body: &ast.LogicalExpression{
			Operator: ast.AndOperator,
			Left: &ast.BinaryExpression{
				Operator: ast.GreaterThanEqualOperator,
				Left: &ast.MemberExpression{
					Object: &ast.Identifier{
						Name: "r",
					},
					Property: &ast.Identifier{
						Name: "_stop",
					},
				},
				Right: &ast.DateTimeLiteral{Value: minStop},
			},
			Right: &ast.BinaryExpression{
				Operator: ast.LessThanEqualOperator,
				Left: &ast.MemberExpression{
					Object: &ast.Identifier{
						Name: "r",
					},
					Property: &ast.Identifier{
						Name: "_start",
					},
				},
				Right: &ast.DateTimeLiteral{Value: maxStart},
			},
		},
	}
}

func columnList(strs ...string) *ast.ArrayExpression {
	list := make([]ast.Expression, len(strs))
	for i, str := range strs {
		if str == model.MetricNameLabel {
			str = "_measurement"
		}
		list[i] = &ast.StringLiteral{Value: str}
	}
	return &ast.ArrayExpression{
		Elements: list,
	}
}

var dropMeasurementCall = call(
	"drop",
	map[string]ast.Expression{
		"columns": &ast.ArrayExpression{
			Elements: []ast.Expression{
				&ast.StringLiteral{Value: "_measurement"},
			},
		},
	},
)

func yieldsFloat(expr promql.Expr) bool {
	switch v := expr.(type) {
	case *promql.NumberLiteral:
		return true
	case *promql.BinaryExpr:
		return yieldsFloat(v.LHS) && yieldsFloat(v.RHS)
	case *promql.UnaryExpr:
		return yieldsFloat(v.Expr)
	case *promql.ParenExpr:
		return yieldsFloat(v.Expr)
	default:
		return false
	}
}

func yieldsTable(expr promql.Expr) bool {
	return !yieldsFloat(expr)
}

func (t *transpiler) transpileExpr(node promql.Node) (ast.Expression, error) {
	switch n := node.(type) {
	case *promql.ParenExpr:
		return t.transpileExpr(n.Expr)
	case *promql.UnaryExpr:
		return t.transpileUnaryExpr(n)
	case *promql.NumberLiteral:
		return &ast.FloatLiteral{Value: n.Val}, nil
	case *promql.StringLiteral:
		return &ast.StringLiteral{Value: n.Val}, nil
	case *promql.VectorSelector:
		return t.transpileInstantVectorSelector(n), nil
	case *promql.MatrixSelector:
		return t.transpileRangeVectorSelector(n), nil
	case *promql.AggregateExpr:
		return t.transpileAggregateExpr(n)
	case *promql.BinaryExpr:
		return t.transpileBinaryExpr(n)
	case *promql.Call:
		return t.transpileCall(n)
	default:
		return nil, fmt.Errorf("PromQL node type %T is not supported yet", t)
	}
}

func (t *transpiler) transpile(node promql.Node) (*ast.File, error) {
	fluxNode, err := t.transpileExpr(node)
	if err != nil {
		return nil, fmt.Errorf("error transpiling expression: %s", err)
	}
	return &ast.File{
		Imports: []*ast.ImportDeclaration{
			{Path: &ast.StringLiteral{Value: "promql"}},
		},
		Body: []ast.Statement{
			&ast.ExpressionStatement{
				Expression: buildPipeline(
					fluxNode,
					// The resolution step evaluation timestamp needs to become the output timestamp.
					call("duplicate", map[string]ast.Expression{
						"column": &ast.StringLiteral{Value: "_stop"},
						"as":     &ast.StringLiteral{Value: "_time"},
					}),
				),
			},
		},
	}, nil
}
