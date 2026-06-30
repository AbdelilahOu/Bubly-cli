package app

import "testing"

func TestExtractBitrate(t *testing.T) {
	cases := map[string]int{
		"128 kbps":     128,
		"130kbps":      130,
		"Best quality": 0,
		"Audio":        0,
		"":             0,
	}
	for in, want := range cases {
		if got := extractBitrate(in); got != want {
			t.Errorf("extractBitrate(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestFixedCol(t *testing.T) {
	cases := []struct {
		value string
		width int
		want  string
	}{
		{"abc", 5, "abc  "},
		{"abcdef", 4, "abc…"},
		{"abc", 3, "abc"},
		{"abc", 1, "…"},
		{"abc", 0, "abc"},
	}
	for _, c := range cases {
		if got := fixedCol(c.value, c.width); got != c.want {
			t.Errorf("fixedCol(%q, %d) = %q, want %q", c.value, c.width, got, c.want)
		}
	}
}

func TestParseAudioFormats(t *testing.T) {
	output := `[youtube] Extracting URL
[info] Available formats for abc123:
ID  EXT   RESOLUTION
--------------------------------------
251 webm  audio only      2  3.50MiB 130k  opus
140 m4a   audio only      2  2.00MiB 128k  m4a_dash
140-drc m4a audio only    2  2.00MiB 128k  m4a_dash`

	formats := ParseAudioFormats(output)
	if len(formats) != 2 {
		t.Fatalf("expected 2 formats (drc skipped), got %d: %+v", len(formats), formats)
	}

	if formats[0].ID != "251" {
		t.Errorf("expected highest-bitrate format first, got %q", formats[0].ID)
	}
	if formats[0].Format != "WebM (Opus)" {
		t.Errorf("expected WebM (Opus), got %q", formats[0].Format)
	}
	if formats[1].Format != "M4A (AAC)" {
		t.Errorf("expected M4A (AAC), got %q", formats[1].Format)
	}
}

func TestParseAudioFormatsFallback(t *testing.T) {
	formats := ParseAudioFormats("no audio lines here")
	if len(formats) != 2 || formats[0].ID != "bestaudio" {
		t.Errorf("expected bestaudio/worstaudio fallback, got %+v", formats)
	}
}

func TestParseVideoFormats(t *testing.T) {
	output := `[youtube] Extracting URL
[info] Available formats for abc123:
ID  EXT   RESOLUTION
--------------------------------------
137 mp4   1920x1080   30  1080p  50.00MiB  video only
251 webm  audio only      2  3.50MiB 130k  opus`

	formats := ParseVideoFormats(output)
	if len(formats) != 1 {
		t.Fatalf("expected 1 video format (audio skipped), got %d: %+v", len(formats), formats)
	}
	if formats[0].ID != "137" {
		t.Errorf("expected id 137, got %q", formats[0].ID)
	}
	if formats[0].Resolution != "1920x1080" {
		t.Errorf("expected resolution 1920x1080, got %q", formats[0].Resolution)
	}
	if formats[0].Quality != "1080p Full HD" {
		t.Errorf("expected quality 1080p Full HD, got %q", formats[0].Quality)
	}
}

func TestParseSubtitleLanguages(t *testing.T) {
	output := `[youtube] Extracting URL
[info] Available subtitles for abc123:
Language Name Formats
en       English vtt
es       Spanish vtt
en-orig  English (original) vtt`

	langs := ParseSubtitleLanguages(output)
	if len(langs) != 2 {
		t.Fatalf("expected 2 languages (en-orig skipped), got %d: %+v", len(langs), langs)
	}
	if langs[0].Code != "en" || langs[0].Name != "English" {
		t.Errorf("expected en/English, got %+v", langs[0])
	}
	if langs[1].Code != "es" || langs[1].Name != "Spanish" {
		t.Errorf("expected es/Spanish, got %+v", langs[1])
	}
}
