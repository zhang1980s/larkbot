package model

type Config struct {
	WideScreenMode bool `json:"wide_screen_mode"`
}
type Text struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type TextElement struct {
	Tag  string `json:"tag"`
	Text Text   `json:"text,omitempty"`
}

type Placeholder struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}
type Value struct {
	Key string `json:"key"`
}
type Options struct {
	Text  Text   `json:"text"`
	Value string `json:"value"`
}
type Extra struct {
	Tag           string      `json:"tag"`
	Placeholder   Placeholder `json:"placeholder"`
	Value         Value       `json:"value"`
	InitialOption string      `json:"initial_option"`
	Options       []Options   `json:"options"`
}
type URLVal struct {
	URL        string `json:"url"`
	AndroidURL string `json:"android_url"`
	IosURL     string `json:"ios_url"`
	PcURL      string `json:"pc_url"`
}
type Href struct {
	URLVal URLVal `json:"urlVal"`
}
type Elements struct {
	Tag     string `json:"tag"`
	Text    Text   `json:"text,omitempty"`
	Extra   Extra  `json:"extra,omitempty"`
	Content string `json:"content,omitempty"`
	Href    Href   `json:"href,omitempty"`
}
type Card struct {
	Config   Config     `json:"config"`
	Elements []Elements `json:"elements"`
}
