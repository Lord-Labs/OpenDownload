package server

import (
	"context"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"

	"github.com/Lord-shaban/OpenDownload/services/api/internal/jobs"
	"github.com/Lord-shaban/OpenDownload/services/api/internal/media"
)

type youtubeResolver []netip.Addr

func (ips youtubeResolver) LookupNetIP(context.Context, string, string) ([]netip.Addr, error) {
	return ips, nil
}

func TestPublicYouTubeReachesEngineAndKeepsURLPolicy(t *testing.T) {
	for _, fixture := range []bool{false, true} {
		s := New(Config{Origin: "https://download.example", Fixture: fixture}, &jobs.Manager{Engine: media.Fixture{}})
		s.Policy.Resolver = youtubeResolver{netip.MustParseAddr("142.250.184.206")}
		for _, source := range []string{"https://youtu.be/2cUkUbB3Gu4", "https://www.youtube.com/watch?v=public", "https://www.youtube-nocookie.com/embed/public"} {
			r := httptest.NewRequest("POST", "/api/v1/analyze", strings.NewReader(`{"url":"`+source+`"}`))
			r.Header.Set("Origin", s.Config.Origin)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Handler().ServeHTTP(w, r)
			if w.Code != 200 || !strings.Contains(w.Body.String(), `"video-1080"`) {
				t.Fatalf("public YouTube did not reach the engine: %d %s", w.Code, w.Body.String())
			}
		}
		if !fixture {
			s.Policy.Resolver = youtubeResolver{netip.MustParseAddr("127.0.0.1")}
			r := httptest.NewRequest("POST", "/api/v1/analyze", strings.NewReader(`{"url":"https://www.youtube.com/watch?v=public"}`))
			r.Header.Set("Origin", s.Config.Origin)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.Handler().ServeHTTP(w, r)
			if w.Code != 400 || !strings.Contains(w.Body.String(), `"unsafe_url"`) {
				t.Fatalf("YouTube bypassed DNS validation: %d %s", w.Code, w.Body.String())
			}
		}
	}
}
