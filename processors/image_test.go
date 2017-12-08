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

func TestCanonicalUrl(t *testing.T) {
	base := "https://test.com"
	base2 := "https://test.com/"
	u1 := "/u1"
	u2 := "www.u2.com"
	u3 := "https://u3.com"
	u4 := "u4/identifier/others"

	r1 := canonicalUrl(base, u1)
	if r1 != "https://test.com/u1" {
		t.Error("Image processor: Canonical URL parsing wrong" + u1)
	}

	r2 := canonicalUrl(base, u2)
	if r2 != "www.u2.com" {
		t.Error("Image processor: Canonical URL parsing wrong " + u2)
	}

	r3 := canonicalUrl(base, u3)
	if r3 != "https://u3.com" {
		t.Error("Image processor: Canonical URL parsing wrong " + u3)
	}

	r4 := canonicalUrl(base, u4)
	if r4 != "https://test.com/u4/identifier/others" {
		t.Error("Image processor: Canonical URL parsing wrong " + u4)
	}

	r5 := canonicalUrl(base2, u1)
	if r5 != "https://test.com/u1" {
		t.Error("Image processor: Canonical URL parsing wrong " + u1)
	}

}
