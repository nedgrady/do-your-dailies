package chorecompletions

import "errors"

// ErrUnknownChoreCompletor means at least one requested completor id does not
// exist on the acting user's account. Surfaced to the client as a 422.
var ErrUnknownChoreCompletor = errors.New("chorecompletions: unknown chore completor for this account")

type CreateRequest struct {
	UserID            uint
	ChoreID           uint
	ChoreCompletorIDs []uint
}
