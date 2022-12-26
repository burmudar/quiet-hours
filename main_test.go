package main

import (
	"os"
	"strings"
	"testing"
)

func TestQuietResponse(t *testing.T) {
	t.Run("marshal and unmarshall return the same thing", func(t *testing.T) {
		resp := QuietResponse{
			IsQuietTime: true,
			WakeUpHour:  7,
			Whoru:       "testing",
		}

		data, err := resp.Marshal()
		if err != nil {
			t.Fatalf("failed to marshal: %s", err)
		}

		got := &QuietResponse{}
		got.Unmarshal(data)

		if resp.Preamble.Version != got.Preamble.Version {
			t.Errorf("version mismatched, got %d wanted %d", got.Preamble.Version, resp.Preamble.Version)
		}

		if resp.IsQuietTime != got.IsQuietTime {
			t.Errorf("got %v wanted %v", got.IsQuietTime, resp.IsQuietTime)
		}

		if resp.WakeUpHour != got.WakeUpHour {
			t.Errorf("got %d wanted %d", got.WakeUpHour, resp.WakeUpHour)
		}
		if resp.Whoru != got.Whoru {
			t.Errorf("got %s wanted %s", got.Whoru, resp.Whoru)
		}
	})
}

func TestQuietQueryEncoding(t *testing.T) {

	t.Run("mashal and unmarshal return the same thing", func(t *testing.T) {
		p := Preamble{
			Version:    1,
			PacketType: QuietQueryType,
		}
		q := QuietQuery{
			Preamble: &p,
			Whoami:   "william",
		}

		data, err := q.Marshal()
		if err != nil {
			t.Fatalf("marshalling failed: %s", err)
		}

		o := &QuietQuery{}
		if _, err := o.Unmarshal(data); err != nil {
			t.Fatalf("unmarshalling failed: %s", err)
		}

		if q.Preamble.Version != o.Preamble.Version {
			t.Errorf("version mismatched, got %d wanted %d", o.Preamble.Version, q.Preamble.Version)
		}
		if q.Whoami != o.Whoami {
			t.Errorf("whoami mismatched, got %s wanted %s", o.Whoami, q.Whoami)
		}
	})
	t.Run("fail when unsupported version", func(t *testing.T) {
		q := QuietQuery{
			Preamble: &Preamble{
				Version:    1,
				PacketType: QuietQueryType,
			},
			Whoami: "william",
		}

		data, err := q.Marshal()
		if err != nil {
			t.Fatalf("marshalling failed: %s", err)
		}

		// First bit is the version
		data[0] = 0

		o := &QuietQuery{}
		if _, err := o.Unmarshal(data); err == nil {
			t.Errorf("expected error during unmarshalling got nil")
		} else if !strings.Contains(err.Error(), "unsupported version") {
			t.Errorf("expected unsupported version error got %q", err.Error())
		}
	})

}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
