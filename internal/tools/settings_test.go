package tools

import "testing"

func TestBuildSettingsBody(t *testing.T) {
	tests := []struct {
		name       string
		args       map[string]any
		wantPrefix string
		wantError  bool
	}{
		{
			name:       "normalizes prefix from JSON",
			args:       map[string]any{"settings_json": `{"task_key_prefix":" ab1 "}`},
			wantPrefix: "AB1",
		},
		{
			name: "typed setting overrides JSON",
			args: map[string]any{
				"settings_json":   `{"task_key_prefix":"OLD","reminder_enabled":false}`,
				"task_key_prefix": "NEW2",
			},
			wantPrefix: "NEW2",
		},
		{
			name:      "rejects invalid prefix",
			args:      map[string]any{"task_key_prefix": "1BAD"},
			wantError: true,
		},
		{
			name:      "rejects non-string prefix from JSON",
			args:      map[string]any{"settings_json": `{"task_key_prefix":123}`},
			wantError: true,
		},
		{
			name:      "requires at least one setting",
			args:      map[string]any{},
			wantError: true,
		},
		{
			name:      "rejects JSON null",
			args:      map[string]any{"settings_json": "null"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := buildSettingsBody(tt.args)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("buildSettingsBody() error = %v", err)
			}
			if got := body["task_key_prefix"]; got != tt.wantPrefix {
				t.Fatalf("task_key_prefix = %v, want %q", got, tt.wantPrefix)
			}
		})
	}
}
