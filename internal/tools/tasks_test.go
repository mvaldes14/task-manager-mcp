package tools

import "testing"

func TestValidateCSVEnum(t *testing.T) {
	allowed := []string{"todo", "doing", "blocked", "done"}
	tests := []struct {
		name      string
		value     any
		wantError bool
	}{
		{name: "missing"},
		{name: "single", value: "todo"},
		{name: "multiple", value: "todo,blocked"},
		{name: "trims spaces", value: "todo, done"},
		{name: "invalid", value: "todo,unknown", wantError: true},
		{name: "empty member", value: "todo,", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := map[string]any{}
			if tt.value != nil {
				args["status"] = tt.value
			}
			err := validateCSVEnum(args, "status", allowed...)
			if tt.wantError && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("validateCSVEnum() error = %v", err)
			}
		})
	}
}
