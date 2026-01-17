package github

import (
	"context"
	"testing"
)

func TestNewTokenClient(t *testing.T) {
	ctx := context.Background()
	token := "dummy-token"
	client, err := NewTokenClient(ctx, token)
	if err != nil {
		t.Fatalf("Failed to create token client: %v", err)
	}
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
}
