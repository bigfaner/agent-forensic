package parser

import "testing"

func TestIsReadTool_KnownAliases(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Read", true},
		{"Read1", true},
		{"Write", false},
		{"Edit", false},
		{"Bash", false},
		{"Agent", false},
		{"SubAgent", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsReadTool(tt.name); got != tt.expected {
				t.Errorf("IsReadTool(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestIsEditTool_KnownAliases(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Write", true},
		{"Edit", true},
		{"Read", false},
		{"Bash", false},
		{"Agent", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEditTool(tt.name); got != tt.expected {
				t.Errorf("IsEditTool(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestIsFileTool_CoversReadAndEdit(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Read", true},
		{"Read1", true},
		{"Write", true},
		{"Edit", true},
		{"Bash", false},
		{"Agent", false},
		{"SubAgent", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFileTool(tt.name); got != tt.expected {
				t.Errorf("IsFileTool(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}

func TestIsAgentTool_KnownAliases(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"Agent", true},
		{"SubAgent", true},
		{"Read", false},
		{"Bash", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAgentTool(tt.name); got != tt.expected {
				t.Errorf("IsAgentTool(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
