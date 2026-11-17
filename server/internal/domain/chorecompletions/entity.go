package chorecompletions

import "errors"

// ErrUnknownChoreCompletor means at least one requested completor id does not
// exist on the acting user's account. Surfaced to the client as a 422.
var ErrUnknownChoreCompletor = errors.New("chorecompletions: unknown chore completor for this account")

// ErrCompletionNotFound means no live completion with that id exists on the
// acting user's account. Surfaced to the client as a 404.
var ErrCompletionNotFound = errors.New("chorecompletions: chore completion not found for this account")

type CreateRequest struct {
	UserID            uint
	ChoreID           uint
	ChoreCompletorIDs []uint
}
