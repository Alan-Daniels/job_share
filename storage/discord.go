package storage

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type DiscordEmbed struct {
	Title  string         `json:"title"`
	Url    string         `json:"url"`
	Fields []DiscordField `json:"fields"`
}
type DiscordMessage struct {
	Content string         `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
}
