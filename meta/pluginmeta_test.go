package meta

import "testing"

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestParsingMetaLines(t *testing.T) {
	pm := New()

	err := pm.ParseMetaLine(" * Plugin Name: example ")
	must(t, err)

	if pm.Name != "example" {
		t.Errorf("Wrong Plugin Name: %s", pm.Name)
	}
}
