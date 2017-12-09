package processor

import (
	"testing"
)

func TestHiddenUrls_simple(t *testing.T) {
	s := `<a href="www.onion1.onion">onion1</a><a href="http://onion2.onion/123"></a>
	<a href="http://not_fetched.com">`
	b := []byte(s)

	urls := hiddenUrls(b)
	if len(urls) != 2 {
		t.Fatal("got wrong number of hidden URLs")
	}
	if urls[0] != "www.onion1.onion" {
		t.Error("www.onion1.onion not parsed correctly")
	}
	if urls[1] != "http://onion2.onion/123" {
		t.Error("www.onion1.onion not parsed correctly")
	}

}
