package db

import "testing"

func TestConfigSetAndGet(t *testing.T) {
	setupTestDB(t)

	if err := SetConfig("postback_url", "https://example.com/hook"); err != nil {
		t.Fatalf("SetConfig failed: %v", err)
	}

	val, err := GetConfig("postback_url")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if val != "https://example.com/hook" {
		t.Errorf("expected 'https://example.com/hook', got %q", val)
	}
}

func TestConfigOverwrite(t *testing.T) {
	setupTestDB(t)

	SetConfig("key", "value1")
	SetConfig("key", "value2")

	val, _ := GetConfig("key")
	if val != "value2" {
		t.Errorf("expected 'value2', got %q", val)
	}
}

func TestConfigNotSet(t *testing.T) {
	setupTestDB(t)

	val, err := GetConfig("nonexistent")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string, got %q", val)
	}
}
