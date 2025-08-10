package appstore_reviews_poller

import (
	"io"
	"log"
	"os"
	"testing"
)

// TestMain suppresses logs during all tests
func TestMain(m *testing.M) {
	// Suppress logs during testing
	log.SetOutput(io.Discard)

	// Run tests
	code := m.Run()

	// Restore log output
	log.SetOutput(os.Stderr)

	// Exit with the same code as the tests
	os.Exit(code)
}
