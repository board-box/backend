package game

type Game struct {
	ID          int64  `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Genre       string `json:"genre" db:"genre"`
	Age         string `json:"age" db:"age"`
	Person      string `json:"person" db:"person"`
	AvgTime     string `json:"avg_time" db:"avg_time"`
	Difficulty  string `json:"complexity" db:"difficulty"`
	Image       string `json:"image" db:"image"`
	Rules       string `json:"rules" db:"rules"`
}
