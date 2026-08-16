package main

import "testing"

func TestResolveDSN(t *testing.T) {
	t.Run("returns error when unset", func(t *testing.T) {
		_, err := resolveDSN(func(string) string { return "" })
		if err == nil {
			t.Fatal("expected an error when DATABASE_URL is unset")
		}
	})

	t.Run("returns DATABASE_URL when set", func(t *testing.T) {
		dsn, err := resolveDSN(func(key string) string {
			if key == "DATABASE_URL" {
				return "host=example dbname=dailies"
			}
			return ""
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if dsn != "host=example dbname=dailies" {
			t.Fatalf("unexpected dsn: %q", dsn)
		}
	})
}

func TestResolvePort(t *testing.T) {
	t.Run("defaults to 8080 when unset", func(t *testing.T) {
		port := resolvePort(func(string) string { return "" })
		if port != "8080" {
			t.Fatalf("expected default port 8080, got %q", port)
		}
	})

	t.Run("respects PORT when set", func(t *testing.T) {
		port := resolvePort(func(key string) string {
			if key == "PORT" {
				return "9090"
			}
			return ""
		})
		if port != "9090" {
			t.Fatalf("expected port 9090, got %q", port)
		}
	})
}
