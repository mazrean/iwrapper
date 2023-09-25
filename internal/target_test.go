package iwrapper

import (
	"os"
	"testing"
)

func TestParseTarget(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description     string
		target          string
		expectedResults []*ParseResult
	}{{
		description: "funcNameなしでも正しくパースできる",
		target:      "normal.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "Normal",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
		}},
	}, {
		description: "funcNameありでも正しくパースできる",
		target:      "func_name.go",
		expectedResults: []*ParseResult{{
			FuncName:   "FuncNameWrapFunc",
			StructName: "FuncName",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
		}},
	}, {
		description:     "targetなしでも正しくパースできる",
		target:          "not_targeted.go",
		expectedResults: nil,
	}, {
		description: "別fileで宣言しても正しくパースできる",
		target:      "other_file_declare.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "OtherFileDeclare",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg:  nil,
				name: "Hijacker",
			}},
		}},
	}, {
		description: "requireが複数あっても正しくパースできる",
		target:      "multi_require.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "MultiRequire",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}, {
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Flusher",
			}},
		}},
	}, {
		description: "optionalが複数あっても正しくパースできる",
		target:      "multi_optional.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "MultiOptional",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}, {
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Flusher",
			}},
		}},
	}, {
		description: "targetが複数あっても正しくパースできる",
		target:      "multi_target.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "MultiTarget1",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
		}, {
			FuncName:   "",
			StructName: "MultiTarget2",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Flusher",
			}},
		}},
	}, {
		description: "括弧内に型定義が1つあっても正しくパースできる",
		target:      "type_in_bracket_outside_comment.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "TypeInBracketOutsideComment",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
		}},
	}, {
		description: "括弧内に型定義があってコメントが括弧内でも正しくパースできる",
		target:      "type_in_bracket_inside_comment.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "TypeInBracketInsideComment",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
		}},
	}, {
		description: "targetが複数括弧内にあっても正しくパースできる",
		target:      "multi_target_in_bracket.go",
		expectedResults: []*ParseResult{{
			FuncName:   "",
			StructName: "MultiTargetInBracket1",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Hijacker",
			}},
		}, {
			FuncName:   "",
			StructName: "MultiTargetInBracket2",
			RequiredInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "ResponseWriter",
			}},
			OptionalInterfaces: []*Interface{{
				pkg: &Package{
					name: "http",
					path: "net/http",
				},
				name: "Flusher",
			}},
		}},
	}}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.description, func(t *testing.T) {
			//t.Parallel()

			target, err := os.Open("testdata/" + testCase.target)
			if err != nil {
				t.Fatal(err)
			}

			pkgName, results, err := ParseTarget(target)
			if err != nil {
				t.Error(err)
			}

			if pkgName != "testdata" {
				t.Errorf("pkgName: expected %q, got %q", "testdata", pkgName)
			}

			if len(results) != len(testCase.expectedResults) {
				t.Errorf("results: expected %d, got %d", len(testCase.expectedResults), len(results))
			}

			if diff := diff(results, testCase.expectedResults); diff != "" {
				t.Errorf("results diff: %s", diff)
			}
		})
	}
}
