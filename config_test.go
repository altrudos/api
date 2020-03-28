package charityhonor

import "testing"

func TestConfig(t *testing.T) {

	c, err := ParseConfig("./config.example.toml")
	if err != nil {
		t.Fatal(err)
	}

	if c.JustGiving.Mode != "staging" {
		t.Fatal("Expecting mode to be staging")
	}
	if c.JustGiving.AppId != "some-app-id" {
		t.Fatal("Expecting app id to be some-app-id")
	}

}
