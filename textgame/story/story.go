package story

var (
	EmptyArc    = Arc{}
	EmptyOption = Option{}
)

// Arc is a story scenario with possible options to take
type Arc struct {
	Name    string
	Title   string   `json:"title"`
	Stories []string `json:"story"`
	Options []Option `json:"options"`
}

// Option groups the text prompt and the arc name to follow
type Option struct {
	Text    string `json:"text"`
	ArcName string `json:"arc"`
}
