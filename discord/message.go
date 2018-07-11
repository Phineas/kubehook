package discord

type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type Image struct {
	URL string `json:"url"`
}

type Thumbnail struct {
	URL string `json:"url"`
}

type Author struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Timestamp   string    `json:"timestamp"`
	Color       int       `json:"color"`
	Footer      Footer    `json:"footer"`
	Image       Image     `json:"image"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Author      Author    `json:"author"`
	Fields      []Field   `json:"fields"`
}

type Message struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds"`
}
