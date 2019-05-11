// DO NOT EDIT: This file is autogenerated via the builtin command.

package promql

import (
	flux "github.com/influxdata/flux"
	ast "github.com/influxdata/flux/ast"
)

func init() {
	flux.RegisterPackage(pkgAST)
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
					Column: 2,
					Line:   27,
				},
				File:   "promql.flux",
				Source: "package promql\n\nbuiltin changes\nbuiltin dayOfMonth\nbuiltin dayOfWeek\nbuiltin daysInMonth\nbuiltin deriv\nbuiltin emptyTable\nbuiltin extrapolatedRate\nbuiltin hour\nbuiltin instantRate\nbuiltin minute\nbuiltin month\nbuiltin resets\nbuiltin timestamp\nbuiltin year\n\n// hack to simulate an imported promql package\npromql = {\n  dayOfMonth:dayOfMonth,\n  dayOfWeek:dayOfWeek,\n  daysInMonth:daysInMonth,\n  hour:hour,\n  minute:minute,\n  month:month,\n  year:year,\n}",
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
						Column: 16,
						Line:   3,
					},
					File:   "promql.flux",
					Source: "builtin changes",
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
							Column: 16,
							Line:   3,
						},
						File:   "promql.flux",
						Source: "changes",
						Start: ast.Position{
							Column: 9,
							Line:   3,
						},
					},
				},
				Name: "changes",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 19,
						Line:   4,
					},
					File:   "promql.flux",
					Source: "builtin dayOfMonth",
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
							Column: 19,
							Line:   4,
						},
						File:   "promql.flux",
						Source: "dayOfMonth",
						Start: ast.Position{
							Column: 9,
							Line:   4,
						},
					},
				},
				Name: "dayOfMonth",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 18,
						Line:   5,
					},
					File:   "promql.flux",
					Source: "builtin dayOfWeek",
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
							Column: 18,
							Line:   5,
						},
						File:   "promql.flux",
						Source: "dayOfWeek",
						Start: ast.Position{
							Column: 9,
							Line:   5,
						},
					},
				},
				Name: "dayOfWeek",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 20,
						Line:   6,
					},
					File:   "promql.flux",
					Source: "builtin daysInMonth",
					Start: ast.Position{
						Column: 1,
						Line:   6,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 20,
							Line:   6,
						},
						File:   "promql.flux",
						Source: "daysInMonth",
						Start: ast.Position{
							Column: 9,
							Line:   6,
						},
					},
				},
				Name: "daysInMonth",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   7,
					},
					File:   "promql.flux",
					Source: "builtin deriv",
					Start: ast.Position{
						Column: 1,
						Line:   7,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   7,
						},
						File:   "promql.flux",
						Source: "deriv",
						Start: ast.Position{
							Column: 9,
							Line:   7,
						},
					},
				},
				Name: "deriv",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 19,
						Line:   8,
					},
					File:   "promql.flux",
					Source: "builtin emptyTable",
					Start: ast.Position{
						Column: 1,
						Line:   8,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 19,
							Line:   8,
						},
						File:   "promql.flux",
						Source: "emptyTable",
						Start: ast.Position{
							Column: 9,
							Line:   8,
						},
					},
				},
				Name: "emptyTable",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 25,
						Line:   9,
					},
					File:   "promql.flux",
					Source: "builtin extrapolatedRate",
					Start: ast.Position{
						Column: 1,
						Line:   9,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 25,
							Line:   9,
						},
						File:   "promql.flux",
						Source: "extrapolatedRate",
						Start: ast.Position{
							Column: 9,
							Line:   9,
						},
					},
				},
				Name: "extrapolatedRate",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   10,
					},
					File:   "promql.flux",
					Source: "builtin hour",
					Start: ast.Position{
						Column: 1,
						Line:   10,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 13,
							Line:   10,
						},
						File:   "promql.flux",
						Source: "hour",
						Start: ast.Position{
							Column: 9,
							Line:   10,
						},
					},
				},
				Name: "hour",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 20,
						Line:   11,
					},
					File:   "promql.flux",
					Source: "builtin instantRate",
					Start: ast.Position{
						Column: 1,
						Line:   11,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 20,
							Line:   11,
						},
						File:   "promql.flux",
						Source: "instantRate",
						Start: ast.Position{
							Column: 9,
							Line:   11,
						},
					},
				},
				Name: "instantRate",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   12,
					},
					File:   "promql.flux",
					Source: "builtin minute",
					Start: ast.Position{
						Column: 1,
						Line:   12,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 15,
							Line:   12,
						},
						File:   "promql.flux",
						Source: "minute",
						Start: ast.Position{
							Column: 9,
							Line:   12,
						},
					},
				},
				Name: "minute",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 14,
						Line:   13,
					},
					File:   "promql.flux",
					Source: "builtin month",
					Start: ast.Position{
						Column: 1,
						Line:   13,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 14,
							Line:   13,
						},
						File:   "promql.flux",
						Source: "month",
						Start: ast.Position{
							Column: 9,
							Line:   13,
						},
					},
				},
				Name: "month",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   14,
					},
					File:   "promql.flux",
					Source: "builtin resets",
					Start: ast.Position{
						Column: 1,
						Line:   14,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 15,
							Line:   14,
						},
						File:   "promql.flux",
						Source: "resets",
						Start: ast.Position{
							Column: 9,
							Line:   14,
						},
					},
				},
				Name: "resets",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 18,
						Line:   15,
					},
					File:   "promql.flux",
					Source: "builtin timestamp",
					Start: ast.Position{
						Column: 1,
						Line:   15,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 18,
							Line:   15,
						},
						File:   "promql.flux",
						Source: "timestamp",
						Start: ast.Position{
							Column: 9,
							Line:   15,
						},
					},
				},
				Name: "timestamp",
			},
		}, &ast.BuiltinStatement{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 13,
						Line:   16,
					},
					File:   "promql.flux",
					Source: "builtin year",
					Start: ast.Position{
						Column: 1,
						Line:   16,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 13,
							Line:   16,
						},
						File:   "promql.flux",
						Source: "year",
						Start: ast.Position{
							Column: 9,
							Line:   16,
						},
					},
				},
				Name: "year",
			},
		}, &ast.VariableAssignment{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 2,
						Line:   27,
					},
					File:   "promql.flux",
					Source: "promql = {\n  dayOfMonth:dayOfMonth,\n  dayOfWeek:dayOfWeek,\n  daysInMonth:daysInMonth,\n  hour:hour,\n  minute:minute,\n  month:month,\n  year:year,\n}",
					Start: ast.Position{
						Column: 1,
						Line:   19,
					},
				},
			},
			ID: &ast.Identifier{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 7,
							Line:   19,
						},
						File:   "promql.flux",
						Source: "promql",
						Start: ast.Position{
							Column: 1,
							Line:   19,
						},
					},
				},
				Name: "promql",
			},
			Init: &ast.ObjectExpression{
				BaseNode: ast.BaseNode{
					Errors: nil,
					Loc: &ast.SourceLocation{
						End: ast.Position{
							Column: 2,
							Line:   27,
						},
						File:   "promql.flux",
						Source: "{\n  dayOfMonth:dayOfMonth,\n  dayOfWeek:dayOfWeek,\n  daysInMonth:daysInMonth,\n  hour:hour,\n  minute:minute,\n  month:month,\n  year:year,\n}",
						Start: ast.Position{
							Column: 10,
							Line:   19,
						},
					},
				},
				Properties: []*ast.Property{&ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 24,
								Line:   20,
							},
							File:   "promql.flux",
							Source: "dayOfMonth:dayOfMonth",
							Start: ast.Position{
								Column: 3,
								Line:   20,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 13,
									Line:   20,
								},
								File:   "promql.flux",
								Source: "dayOfMonth",
								Start: ast.Position{
									Column: 3,
									Line:   20,
								},
							},
						},
						Name: "dayOfMonth",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 24,
									Line:   20,
								},
								File:   "promql.flux",
								Source: "dayOfMonth",
								Start: ast.Position{
									Column: 14,
									Line:   20,
								},
							},
						},
						Name: "dayOfMonth",
					},
				}, &ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 22,
								Line:   21,
							},
							File:   "promql.flux",
							Source: "dayOfWeek:dayOfWeek",
							Start: ast.Position{
								Column: 3,
								Line:   21,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 12,
									Line:   21,
								},
								File:   "promql.flux",
								Source: "dayOfWeek",
								Start: ast.Position{
									Column: 3,
									Line:   21,
								},
							},
						},
						Name: "dayOfWeek",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 22,
									Line:   21,
								},
								File:   "promql.flux",
								Source: "dayOfWeek",
								Start: ast.Position{
									Column: 13,
									Line:   21,
								},
							},
						},
						Name: "dayOfWeek",
					},
				}, &ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 26,
								Line:   22,
							},
							File:   "promql.flux",
							Source: "daysInMonth:daysInMonth",
							Start: ast.Position{
								Column: 3,
								Line:   22,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 14,
									Line:   22,
								},
								File:   "promql.flux",
								Source: "daysInMonth",
								Start: ast.Position{
									Column: 3,
									Line:   22,
								},
							},
						},
						Name: "daysInMonth",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 26,
									Line:   22,
								},
								File:   "promql.flux",
								Source: "daysInMonth",
								Start: ast.Position{
									Column: 15,
									Line:   22,
								},
							},
						},
						Name: "daysInMonth",
					},
				}, &ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 12,
								Line:   23,
							},
							File:   "promql.flux",
							Source: "hour:hour",
							Start: ast.Position{
								Column: 3,
								Line:   23,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 7,
									Line:   23,
								},
								File:   "promql.flux",
								Source: "hour",
								Start: ast.Position{
									Column: 3,
									Line:   23,
								},
							},
						},
						Name: "hour",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 12,
									Line:   23,
								},
								File:   "promql.flux",
								Source: "hour",
								Start: ast.Position{
									Column: 8,
									Line:   23,
								},
							},
						},
						Name: "hour",
					},
				}, &ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 16,
								Line:   24,
							},
							File:   "promql.flux",
							Source: "minute:minute",
							Start: ast.Position{
								Column: 3,
								Line:   24,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 9,
									Line:   24,
								},
								File:   "promql.flux",
								Source: "minute",
								Start: ast.Position{
									Column: 3,
									Line:   24,
								},
							},
						},
						Name: "minute",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 16,
									Line:   24,
								},
								File:   "promql.flux",
								Source: "minute",
								Start: ast.Position{
									Column: 10,
									Line:   24,
								},
							},
						},
						Name: "minute",
					},
				}, &ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 14,
								Line:   25,
							},
							File:   "promql.flux",
							Source: "month:month",
							Start: ast.Position{
								Column: 3,
								Line:   25,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 8,
									Line:   25,
								},
								File:   "promql.flux",
								Source: "month",
								Start: ast.Position{
									Column: 3,
									Line:   25,
								},
							},
						},
						Name: "month",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 14,
									Line:   25,
								},
								File:   "promql.flux",
								Source: "month",
								Start: ast.Position{
									Column: 9,
									Line:   25,
								},
							},
						},
						Name: "month",
					},
				}, &ast.Property{
					BaseNode: ast.BaseNode{
						Errors: nil,
						Loc: &ast.SourceLocation{
							End: ast.Position{
								Column: 12,
								Line:   26,
							},
							File:   "promql.flux",
							Source: "year:year",
							Start: ast.Position{
								Column: 3,
								Line:   26,
							},
						},
					},
					Key: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 7,
									Line:   26,
								},
								File:   "promql.flux",
								Source: "year",
								Start: ast.Position{
									Column: 3,
									Line:   26,
								},
							},
						},
						Name: "year",
					},
					Value: &ast.Identifier{
						BaseNode: ast.BaseNode{
							Errors: nil,
							Loc: &ast.SourceLocation{
								End: ast.Position{
									Column: 12,
									Line:   26,
								},
								File:   "promql.flux",
								Source: "year",
								Start: ast.Position{
									Column: 8,
									Line:   26,
								},
							},
						},
						Name: "year",
					},
				}},
			},
		}},
		Imports: nil,
		Name:    "promql.flux",
		Package: &ast.PackageClause{
			BaseNode: ast.BaseNode{
				Errors: nil,
				Loc: &ast.SourceLocation{
					End: ast.Position{
						Column: 15,
						Line:   1,
					},
					File:   "promql.flux",
					Source: "package promql",
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
							Column: 15,
							Line:   1,
						},
						File:   "promql.flux",
						Source: "promql",
						Start: ast.Position{
							Column: 9,
							Line:   1,
						},
					},
				},
				Name: "promql",
			},
		},
	}},
	Package: "promql",
	Path:    "promql",
}
