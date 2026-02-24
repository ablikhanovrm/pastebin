package cache

import (
	"fmt"
	"time"
)

func pasteKey(id string) string {
	return "paste:" + id
}

func pasteListKey(limit int32, cursor *time.Time) string {
	if cursor != nil {
		return fmt.Sprintf("paste_list:%s:%d", cursor.Format(time.RFC3339), limit)
	}

	return fmt.Sprintf("paste_list:first:%d", limit)
}

func pasteContentKey(id string) string {
	return "paste_content:" + id
}
