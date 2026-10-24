package media

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"testing"
)

func TestYouTubePublicFormatsAndRestrictions(t *testing.T) {
	data, err := os.ReadFile("testdata/youtube-public.json")
	if err != nil {
		t.Fatal(err)
	}
	var raw RawInfo
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{"https://youtu.be/2cUkUbB3Gu4", "https://example.com/embed"} {
		a, err := Normalize(raw, source)
		if err != nil {
			t.Fatal(err)
		}
		if a.Platform != "Youtube" || a.Options[0].Selector != "137+140" || a.Options[0].Extension != "mp4" {
			t.Fatalf("bad YouTube selection: %+v", a)
		}
		if !slices.ContainsFunc(a.Options, func(o Option) bool { return o.ID == "audio-mp3" && o.Selector == "140" }) {
			t.Fatal("YouTube audio conversion is missing")
		}
		for _, o := range a.Options {
			if !validSelector(o.Selector) {
				t.Fatalf("invalid selector: %+v", o)
			}
		}
	}
	for _, restricted := range []RawInfo{{Extractor: "Youtube", DRM: true}, {Extractor: "Youtube", Live: true}, {Extractor: "Youtube", Availability: "needs_auth"}, {Extractor: "Youtube", Availability: "private"}} {
		if _, err := Normalize(restricted, "https://youtu.be/2cUkUbB3Gu4"); !errors.Is(err, ErrAccess) {
			t.Fatalf("restricted YouTube accepted: %v", err)
		}
	}
	for i := range raw.Formats {
		if raw.Formats[i].ID == "137" {
			raw.Formats[i].Bytes, raw.Formats[i].Approx = 0, 0
		}
	}
	a, err := Normalize(raw, "https://youtu.be/2cUkUbB3Gu4")
	if err != nil || a.Options[0].Bytes != 0 {
		t.Fatalf("unknown video size was presented as a known merged size: %+v %v", a.Options, err)
	}
}
