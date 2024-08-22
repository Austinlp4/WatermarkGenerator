package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
	"watermark-generator/db"
)

func TestMongoConnection(t *testing.T) {
	db.Connect()

	// Test if the client is not nil
	client := db.GetClient()
	if client == nil {
		t.Fatal("Client is nil after Connect")
	}

	// Test ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
}

func redactURI(uri string) string {
	parts := strings.Split(uri, "@")
	if len(parts) > 1 {
		credentialParts := strings.Split(parts[0], "://")
		if len(credentialParts) > 1 {
			return credentialParts[0] + "://<redacted>@" + parts[1]
		}
	}
	return "<unable to redact>"
}
