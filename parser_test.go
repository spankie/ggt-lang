package lang

import (
	"strings"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	template := `f('hello {{ .foo }}', {
		bar: 'foo',
	});`

	parser := NewParser(strings.NewReader(template))
	stmt, err := parser.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if stmt.Template != "'hello {{ .foo }}'" {
		t.Fatalf("unexpected template: %s", stmt.Template)
	}
	if stmt.Vars["bar"] != "'foo'" {
		t.Fatalf("unexpected value for foo: %s", stmt.Vars["bar"])
	}

	result, err := stmt.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if result != "hello foo" {
		t.Fatalf("unexpected result: %s", result)
	}
}
