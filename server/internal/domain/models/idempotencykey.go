package models

import "time"

// IdempotencyKey records one POST request made with an Idempotency-Key header,
// so repeats of the same request replay the original response instead of
// running the handler again. Keys are scoped to the authenticated principal
// that made the request (a user today; could be an API client or tenant later),
// so one principal can never replay or block another principal's requests.
type IdempotencyKey struct {
	ID                  uint   `gorm:"primaryKey"`
	Scope               string `gorm:"size:128;not null;uniqueIndex:idx_idempotency_scope_key"`
	Key                 string `gorm:"size:64;not null;uniqueIndex:idx_idempotency_scope_key"`
	RequestHash         string `gorm:"size:64;not null"`
	State               string `gorm:"size:16;not null"`
	ResponseStatus      int
	ResponseContentType string
	ResponseBody        []byte
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
