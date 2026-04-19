package utctime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSONNormalizesOffsetToUTC(t *testing.T) {
	t.Parallel()
	var value Time

	err := json.Unmarshal([]byte(`"2026-04-07T12:00:00+02:00"`), &value)
	if err != nil {
		t.Fatalf("unmarshal time: %v", err)
	}

	assert.Equal(t, "2026-04-07T10:00:00Z", value.Format(time.RFC3339))
}

func TestMarshalJSONEmitsUTCZulu(t *testing.T) {
	t.Parallel()
	value := Time{Time: time.Date(2026, time.April, 7, 12, 0, 0, 0, time.FixedZone("UTC+2", 2*60*60))}

	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal time: %v", err)
	}

	assert.Equal(t, `"2026-04-07T10:00:00Z"`, string(encoded))
}
