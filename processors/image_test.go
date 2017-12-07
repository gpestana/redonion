package processor

import (
	"testing"
)

func TestImageParser_simple(t *testing.T) {
	s := `<p>Images:</p><ul><li><img src="foo.jpeg"/></li><li><img src="/bar/baz.png"/></li></ul>`
	b := []byte(s)

	imgs := images(b)
	if len(imgs) != 2 {
		t.Fatalf("Got wrong number of parsed image URL")
	}
	if imgs[0] != "foo.jpeg" {
		t.Error("foo.jpeg not parsed correctly")
	}

	if imgs[1] != "/bar/baz.png" {
		t.Error("/bar/baz.png not parsed correctly")
	}

}
