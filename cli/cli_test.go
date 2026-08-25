package cli

import (
	"strings"
	"testing"

	"devicecert/testkit"
)

func TestCLIUsage(t *testing.T) {
	db, manager, err := testkit.OpenFixture("cli")
	if err != nil {
		t.Fatal(err)
	}
	defer testkit.Cleanup(db, "cli")
	app := New(manager)
	if output, err := app.Execute(nil); err != nil || !strings.Contains(output, "register") {
		t.Fatalf("usage failed: %v", err)
	}
	if output, err := app.ExecuteLine("register DC-12345678 ALPHA-1"); err != nil || !strings.Contains(output, "registered") {
		t.Fatalf("register failed: %v", err)
	}
	if output, err := app.Execute([]string{"issue", "DC-12345678"}); err != nil || !strings.Contains(output, "issued") {
		t.Fatalf("issue failed: %v", err)
	}
}
