package collection

type Collection struct {
	ID      int64   `json:"id" db:"id"`
	UserID  int64   `json:"user_id" db:"user_id"`
	Name    string  `json:"name" db:"name"`
	Pinned  bool    `json:"pinned" db:"pinned"`
	GameIDs []int64 `json:"game_ids,omitempty"`
}
