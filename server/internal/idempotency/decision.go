package idempotency

import (
	"crypto/sha256"
	"encoding/hex"

	"do-your-dailies/server/internal/domain/models"
)

const (
	StateInProgress = "in_progress"
	StateCompleted  = "completed"
)

type Action int

const (
	// ActionRun means this request is the first with its key and must run the handler.
	ActionRun Action = iota
	// ActionReplay means the earlier attempt finished; return its stored response.
	ActionReplay
	// ActionInProgress means an earlier attempt with this key is still running (or
	// crashed while running). The client should retry under a new key.
	ActionInProgress
	// ActionMismatch means the key was reused for a different request.
	ActionMismatch
)

// Fingerprint identifies a request for key-reuse detection. Two requests with the
// same key but different fingerprints are a client bug, not a retry.
func Fingerprint(method, path string, body []byte) string {
	sum := sha256.New()
	sum.Write([]byte(method))
	sum.Write([]byte{'\n'})
	sum.Write([]byte(path))
	sum.Write([]byte{'\n'})
	sum.Write(body)
	return hex.EncodeToString(sum.Sum(nil))
}

// Decide chooses what to do with a request, given the row already stored under its
// key. It is pure so the rules can be read and reasoned about without a database.
func Decide(existing models.IdempotencyKey, requestHash string) Action {
	if existing.RequestHash != requestHash {
		return ActionMismatch
	}
	if existing.State == StateCompleted {
		return ActionReplay
	}
	return ActionInProgress
}
