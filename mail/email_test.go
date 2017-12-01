package mail

import (
	"os"
	"testing"

	"github.com/benjamw/golibs/test"

	"github.com/benjamw/gogame/config"
)

var (
	email = "benjamwelker@gmail.com" // set this to a valid email address for the send email test
)

func TestMain(m *testing.M) {
	test.InitCtx()

	config.GenRoot()
	config.RootURL = "UNIT_TESTING"

	runVal := m.Run()
	test.ReleaseCtx()
	os.Exit(runVal)
}

func TestParseTemplate(t *testing.T) {
	params := make(map[string]interface{}, 0)
	split, err := parseTemplates("welcome", params)
	if err != nil {
		t.Fatalf("mail.parseTemplates threw an error: %v", err)
	}
	if len(split) != 3 {
		t.Fatalf("mail.parseTemplates did not return the correct number of items. Wanted: 3; Got: %d", len(split))
	}
}

func TestFromTemplate(t *testing.T) {
	ctx := test.GetCtx()

	err := FromTemplate(ctx, "welcome", []string{email}, nil)
	if err != nil {
		t.Fatalf("mail.FromTemplate threw an error: %v", err)
	}
}
