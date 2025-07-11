package workflow

const (
	AVAILABLE int64 = iota
	RENTED
)

type UserProfile struct {
	ID       string
	Username string
	Email    string
	Items    []string
	Address  string
}

type Review struct {
	ID        string
	LuggageID string
	Text      string
	Rating    int64
	User      string
}

type LuggageItem struct {
	ID      string
	Color   string
	Length  int64
	Breadth int64
	Height  int64
	Price   float64
}

type LuggageInfo struct {
	Item    LuggageItem
	Reviews []Review
}

type Reservation struct {
	LuggageID string
	StartDate string
	EndDate   string
	User      string
}
