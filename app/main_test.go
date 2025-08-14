package main

import (
	"testing"
)

func TestExtractSSHPort(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0.0.0.0:2222->22/tcp", "2222"},
		{"0.0.0.0:2223->22/tcp", "2223"},
		{"0.0.0.0:2224->22/tcp", "2224"},
		{"invalid", "N/A"},
		{"", "N/A"},
	}

	for _, test := range tests {
		result := extractSSHPort(test.input)
		if result != test.expected {
			t.Errorf("extractSSHPort(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestExtractHostname(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"lab-01", "lab-01"},
		{"lab-02", "lab-02"},
		{"lab-10", "lab-10"},
		{"lab-99", "lab-99"},
		{"invalid", "unknown"},
		{"", "unknown"},
	}

	for _, test := range tests {
		result := extractHostname(test.input)
		if result != test.expected {
			t.Errorf("extractHostname(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestGenerateDockerCompose(t *testing.T) {
	tests := []struct {
		containerCount int
		expectError    bool
	}{
		{2, false},
		{5, false},
		{1, false},
		{10, false},
	}

	for _, test := range tests {
		err := generateDockerCompose(test.containerCount)
		if test.expectError && err == nil {
			t.Errorf("generateDockerCompose(%d) expected error, got nil", test.containerCount)
		}
		if !test.expectError && err != nil {
			t.Errorf("generateDockerCompose(%d) unexpected error: %v", test.containerCount, err)
		}
	}
}
