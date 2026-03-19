package main

import "testing"

func TestResolveParent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		projectID string
		parent    string
		want      string
		wantErr   bool
	}{
		{
			name:      "uses explicit parent",
			projectID: "ignored-project",
			parent:    "folders/123456789",
			want:      "folders/123456789",
		},
		{
			name:      "builds parent from project id",
			projectID: "demo-project",
			want:      "projects/demo-project",
		},
		{
			name:    "errors without either input",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resolveParent(tt.projectID, tt.parent)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveParent returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("resolveParent = %q, want %q", got, tt.want)
			}
		})
	}
}
