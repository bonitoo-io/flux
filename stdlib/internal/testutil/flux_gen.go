// DO NOT EDIT: This file is autogenerated via the builtin command.

package testutil

import (
	ast "github.com/influxdata/flux/ast"
	runtime "github.com/influxdata/flux/runtime"
)

func init() {
	runtime.RegisterPackage(pkgAST)
}

var pkgAST = &ast.Package{
	BaseNode: ast.BaseNode{
		Errors: nil,
		Loc:    nil,
	},
	Files: []*ast.File{&ast.File{
		BaseNode: ast.BaseNode{
			Errors: nil,
			Loc: &ast.SourceLocation{
				End: ast.Position{
					Column: 19,
					Line:   5,
				},
				File:   "testutil.flux",
				Source: "package testutil\n\nbuiltin fail : () => bool\nbuiltin yield : (<-v: A) => A\nbuiltin makeRecord",
				Start: ast.Position{
					Column: 1,
					Line:   1,
				},
			},
		},
		Body: []ast.Statement{&ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   3,
					},
					File:   "testutil.flux",
					Source: "builtin fail",
					Start: ast.Position{
						Column: 1,
						Line:   3,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 13,
							Line:   3,
						},
						File:   "testutil.flux",
						Source: "fail",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "fail",
			},
			Ty: ast.TypeExpression{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 26,
							Line:   3,
						},
						File:   "testutil.flux",
						Source: "() => bool",
						Start: ast.Position{
							Column: 16,
							Line:   3,
						},
					},
				},
				Constraints: []*ast.TypeConstraint{},
				Ty: &ast.FunctionType{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 26,
								Line:   3,
							},
							File:   "testutil.flux",
							Source: "() => bool",
							Start: ast.Position{
								Column: 16,
								Line:   3,
							},
						},
					},
					Parameters: []*ast.ParameterType{},
					Return: &ast.NamedType{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 26,
									Line:   3,
								},
								File:   "testutil.flux",
								Source: "bool",
								Start: ast.Position{
									Column: 22,
									Line:   3,
								},
							},
						},
						ID: &ast.Identifier{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 26,
										Line:   3,
									},
									File:   "testutil.flux",
									Source: "bool",
									Start: ast.Position{
										Column: 22,
										Line:   3,
									},
								},
							},
							Name: "bool",
						},
					},
				},
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   4,
					},
					File:   "testutil.flux",
					Source: "builtin yield",
					Start: ast.Position{
						Column: 1,
						Line:   4,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   4,
						},
						File:   "testutil.flux",
						Source: "yield",
						Start: ast.Position{
							Column: 9,
							Line:   4,
						},
					},
				},
				Name: "yield",
			},
			Ty: ast.TypeExpression{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 30,
							Line:   4,
						},
						File:   "testutil.flux",
						Source: "(<-v: A) => A",
						Start: ast.Position{
							Column: 17,
							Line:   4,
						},
					},
				},
				Constraints: []*ast.TypeConstraint{},
				Ty: &ast.FunctionType{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 30,
								Line:   4,
							},
							File:   "testutil.flux",
							Source: "(<-v: A) => A",
							Start: ast.Position{
								Column: 17,
								Line:   4,
							},
						},
					},
					Parameters: []*ast.ParameterType{&ast.ParameterType{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 24,
									Line:   4,
								},
								File:   "testutil.flux",
								Source: "<-v: A",
								Start: ast.Position{
									Column: 18,
									Line:   4,
								},
							},
						},
						Kind: "Pipe",
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 21,
										Line:   4,
									},
									File:   "testutil.flux",
									Source: "v",
									Start: ast.Position{
										Column: 20,
										Line:   4,
									},
								},
							},
							Name: "v",
						},
						Ty: &ast.TvarType{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 24,
										Line:   4,
									},
									File:   "testutil.flux",
									Source: "A",
									Start: ast.Position{
										Column: 23,
										Line:   4,
									},
								},
							},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{
									Errors: nil,
									Loc: &ast.SourceLocation{
										End: ast.Position{
											Column: 24,
											Line:   4,
										},
										File:   "testutil.flux",
										Source: "A",
										Start: ast.Position{
											Column: 23,
											Line:   4,
										},
									},
								},
								Name: "A",
							},
						},
					}},
					Return: &ast.TvarType{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 30,
									Line:   4,
								},
								File:   "testutil.flux",
								Source: "A",
								Start: ast.Position{
									Column: 29,
									Line:   4,
								},
							},
						},
						ID: &ast.Identifier{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 30,
										Line:   4,
									},
									File:   "testutil.flux",
									Source: "A",
									Start: ast.Position{
										Column: 29,
										Line:   4,
									},
								},
							},
							Name: "A",
						},
					},
				},
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 19,
						Line:   5,
					},
					File:   "testutil.flux",
					Source: "builtin makeRecord",
					Start: ast.Position{
						Column: 1,
						Line:   5,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 19,
							Line:   5,
						},
						File:   "testutil.flux",
						Source: "makeRecord",
						Start: ast.Position{
							Column: 9,
							Line:   5,
						},
					},
				},
				Name: "makeRecord",
			},
			Ty: ast.TypeExpression{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 60,
							Line:   5,
						},
						File:   "testutil.flux",
						Source: "(o: A) => B where A: Record, B: Record",
						Start: ast.Position{
							Column: 22,
							Line:   5,
						},
					},
				},
				Constraints: []*ast.TypeConstraint{&ast.TypeConstraint{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 49,
								Line:   5,
							},
							File:   "testutil.flux",
							Source: "A: Record",
							Start: ast.Position{
								Column: 40,
								Line:   5,
							},
						},
					},
					Kinds: []*ast.Identifier{&ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 49,
									Line:   5,
								},
								File:   "testutil.flux",
								Source: "Record",
								Start: ast.Position{
									Column: 43,
									Line:   5,
								},
							},
						},
						Name: "Record",
					}},
					Tvar: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 41,
									Line:   5,
								},
								File:   "testutil.flux",
								Source: "A",
								Start: ast.Position{
									Column: 40,
									Line:   5,
								},
							},
						},
						Name: "A",
					},
				}, &ast.TypeConstraint{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 60,
								Line:   5,
							},
							File:   "testutil.flux",
							Source: "B: Record",
							Start: ast.Position{
								Column: 51,
								Line:   5,
							},
						},
					},
					Kinds: []*ast.Identifier{&ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 60,
									Line:   5,
								},
								File:   "testutil.flux",
								Source: "Record",
								Start: ast.Position{
									Column: 54,
									Line:   5,
								},
							},
						},
						Name: "Record",
					}},
					Tvar: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 52,
									Line:   5,
								},
								File:   "testutil.flux",
								Source: "B",
								Start: ast.Position{
									Column: 51,
									Line:   5,
								},
							},
						},
						Name: "B",
					},
				}},
				Ty: &ast.FunctionType{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 33,
								Line:   5,
							},
							File:   "testutil.flux",
							Source: "(o: A) => B",
							Start: ast.Position{
								Column: 22,
								Line:   5,
							},
						},
					},
					Parameters: []*ast.ParameterType{&ast.ParameterType{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 27,
									Line:   5,
								},
								File:   "testutil.flux",
								Source: "o: A",
								Start: ast.Position{
									Column: 23,
									Line:   5,
								},
							},
						},
						Kind: "Required",
						Name: &ast.Identifier{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 24,
										Line:   5,
									},
									File:   "testutil.flux",
									Source: "o",
									Start: ast.Position{
										Column: 23,
										Line:   5,
									},
								},
							},
							Name: "o",
						},
						Ty: &ast.TvarType{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 27,
										Line:   5,
									},
									File:   "testutil.flux",
									Source: "A",
									Start: ast.Position{
										Column: 26,
										Line:   5,
									},
								},
							},
							ID: &ast.Identifier{
								BaseNode: ast.BaseNode{
									Errors: nil,
									Loc: &ast.SourceLocation{
										End: ast.Position{
											Column: 27,
											Line:   5,
										},
										File:   "testutil.flux",
										Source: "A",
										Start: ast.Position{
											Column: 26,
											Line:   5,
										},
									},
								},
								Name: "A",
							},
						},
					}},
					Return: &ast.TvarType{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 33,
									Line:   5,
								},
								File:   "testutil.flux",
								Source: "B",
								Start: ast.Position{
									Column: 32,
									Line:   5,
								},
							},
						},
						ID: &ast.Identifier{
							BaseNode: ast.BaseNode{
								Errors: nil,
								Loc: &ast.SourceLocation{
									End: ast.Position{
										Column: 33,
										Line:   5,
									},
									File:   "testutil.flux",
									Source: "B",
									Start: ast.Position{
										Column: 32,
										Line:   5,
									},
								},
							},
							Name: "B",
						},
					},
				},
			},
		}},
		Imports:  nil,
		Metadata: "parser-type=rust",
		Name:     "testutil.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 17,
						Line:   1,
					},
					File:   "testutil.flux",
					Source: "package testutil",
					Start: ast.Position{
						Column: 1,
						Line:   1,
					},
				},
			},
			Name: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 17,
							Line:   1,
						},
						File:   "testutil.flux",
						Source: "testutil",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "testutil",
			},
		},
	}},
	Package: "testutil",
	Path:    "internal/testutil",
}
