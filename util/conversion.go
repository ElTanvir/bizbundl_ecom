package util

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", u.Bytes[0:4], u.Bytes[4:6], u.Bytes[6:8], u.Bytes[8:10], u.Bytes[10:16])
}

// Better yet, just use uuid.UUID if we can convert.
// pgtype.UUID Has specific byte access.
// Let's use a simpler approach if we imported google/uuid:
// But we are using pgtype.UUID.
// Actually pgtype.UUID.Bytes is [16]byte.

func Int32ToString(n int32) string {
	return strconv.Itoa(int(n))
}
