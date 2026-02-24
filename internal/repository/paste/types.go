package paste

import "time"

type GetPastesAfterCursorParams struct {
	UserId int64
	Limit  int32
	Cursor time.Time
}

type GetPastesFirstPageParams struct {
	UserId int64
	Limit  int32
}

type GetMyPastesFirstPageParams struct {
	UserId int64
	Limit  int32
}

type GetMyPastesAfterCursorParams struct {
	UserId int64
	Limit  int32
	Cursor time.Time
}
