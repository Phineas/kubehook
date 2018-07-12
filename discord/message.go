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
	Title       string    `json:"title",omitempty`
	Description string    `json:"description",omitempty`
	URL         string    `json:"url",omitempty`
	Timestamp   string    `json:"timestamp",omitempty`
	Color       int       `json:"color",omitempty`
	Footer      Footer    `json:"footer",omitempty`
	Image       Image     `json:"image",omitempty`
	Thumbnail   Thumbnail `json:"thumbnail",omitempty`
	Author      Author    `json:"author",omitempty`
	Fields      []Field   `json:"fields",omitempty`
}

type Message struct {
	Content string  `json:"content"`
	Embeds  []Embed `json:"embeds",omitempty`
}
