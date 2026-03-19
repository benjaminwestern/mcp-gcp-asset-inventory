package main

import (
	"testing"
	"time"

	"cloud.google.com/go/asset/apiv1/assetpb"
)

func TestParseContentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		raw         string
		defaultType assetpb.ContentType
		want        assetpb.ContentType
		wantErr     bool
	}{
		{
			name:        "falls back to default",
			raw:         "",
			defaultType: assetpb.ContentType_RESOURCE,
			want:        assetpb.ContentType_RESOURCE,
		},
		{
			name:        "accepts lowercase alias",
			raw:         "iam_policy",
			defaultType: assetpb.ContentType_RESOURCE,
			want:        assetpb.ContentType_IAM_POLICY,
		},
		{
			name:        "accepts dashed alias",
			raw:         "access-policy",
			defaultType: assetpb.ContentType_RESOURCE,
			want:        assetpb.ContentType_ACCESS_POLICY,
		},
		{
			name:        "rejects invalid value",
			raw:         "not-a-real-type",
			defaultType: assetpb.ContentType_RESOURCE,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseContentType(tt.raw, tt.defaultType)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("parseContentType returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("parseContentType = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTimeWindow(t *testing.T) {
	t.Parallel()

	t.Run("requires start time when requested", func(t *testing.T) {
		t.Parallel()

		_, err := parseTimeWindow("", "2026-03-19T10:30:00Z", true)
		if err == nil {
			t.Fatalf("expected an error but got none")
		}
	})

	t.Run("rejects reversed windows", func(t *testing.T) {
		t.Parallel()

		_, err := parseTimeWindow(
			"2026-03-20T10:30:00Z",
			"2026-03-19T10:30:00Z",
			false,
		)
		if err == nil {
			t.Fatalf("expected an error but got none")
		}
	})

	t.Run("builds a valid window", func(t *testing.T) {
		t.Parallel()

		got, err := parseTimeWindow(
			"2026-03-19T10:30:00Z",
			"2026-03-20T10:30:00Z",
			true,
		)
		if err != nil {
			t.Fatalf("parseTimeWindow returned error: %v", err)
		}
		if got == nil {
			t.Fatalf("expected a time window but got nil")
		}

		wantStart := time.Date(2026, 3, 19, 10, 30, 0, 0, time.UTC)
		wantEnd := time.Date(2026, 3, 20, 10, 30, 0, 0, time.UTC)

		if !got.GetStartTime().AsTime().Equal(wantStart) {
			t.Fatalf("start_time = %v, want %v", got.GetStartTime().AsTime(), wantStart)
		}
		if !got.GetEndTime().AsTime().Equal(wantEnd) {
			t.Fatalf("end_time = %v, want %v", got.GetEndTime().AsTime(), wantEnd)
		}
	})
}

func TestNormalizePageSize(t *testing.T) {
	t.Parallel()

	if got := normalizePageSize(0, 100, 500); got != 100 {
		t.Fatalf("normalizePageSize default = %d, want 100", got)
	}

	if got := normalizePageSize(900, 100, 500); got != 500 {
		t.Fatalf("normalizePageSize clamp = %d, want 500", got)
	}

	if got := normalizePageSize(200, 100, 500); got != 200 {
		t.Fatalf("normalizePageSize exact = %d, want 200", got)
	}
}
