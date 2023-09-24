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
			FuncName: "NormalWrapper",
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
			FuncName: "FuncNameWrapFunc",
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
			FuncName: "OtherFileDeclareWrapper",
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
			FuncName: "MultiRequireWrapper",
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
			FuncName: "MultiOptionalWrapper",
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
			FuncName: "MultiTarget1Wrapper",
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
			FuncName: "MultiTarget2Wrapper",
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
			t.Parallel()

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
