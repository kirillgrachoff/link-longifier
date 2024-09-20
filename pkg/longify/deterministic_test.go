package longify_test

import (
	"testing"
	"unicode/utf8"

	"github.com/kirillgrachoff/link-longifier/pkg/longify"
	"pgregory.net/rapid"
)

func TestJustWorks(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		link := rapid.StringMatching(`https?\:\/\/[a-z0-9_]{4,}\.[a-z]{3,9}(\/[a-zA-Z0-9_])*\/?`).Draw(t, "link")
		longLink, err := longify.ForwardDeterminisic(link)
		if len(longLink) < len(link) {
			t.Errorf("link %v is shortened to %v", link, longLink)
			return
		}
		if err != nil {
			t.Error(err)
			return
		}
		if !utf8.Valid([]byte(longLink)) {
			t.Errorf("string %v is not a valid utf8 string", []byte(longLink))
			return
		}
		t.Logf("link %v is longified to %v", link, longLink)
		shortenLink, err := longify.BackwardDeterministic(longLink)
		if err != nil {
			t.Error(err)
			return
		}
		if link != shortenLink {
			t.Errorf("link %+v and link %+v are not equal", link, shortenLink)
		}
	})
}

func BenchmarkForwardBackward(b *testing.B) {
	link := `https://example.com/example-link-to-some_resource`
	for i := 0; i < b.N; i++ {
		longLink, _ := longify.ForwardDeterminisic(link)
		shortLink, _ := longify.BackwardDeterministic(longLink)
		_ = shortLink
	}
}
