package models

type Article struct {
	URL      string `toml:"url"`
	Title    string
	Type     string `toml:"type"`
	Category string `toml:"cat"`
	Summary  string `toml:"summary"`
}

type NewsletterData struct {
	Title       string
	IssueNo     int
	PubDate     string
	WelcomeText string
	Articles    []Article
	Images      []Image
}

type SourceFile struct {
	Metadata SourceFileMetadata `toml:"metadata"`
	Articles []Article          `toml:"articles"`
	Images   []Image            `toml:"images"`
}

type SourceFileMetadata struct {
	Title       string `toml:"title"`
	IssueNo     int    `toml:"no"`
	PubDate     string `toml:"date"`
	WelcomeText string `toml:"welcome"`
}

type Image struct {
	FilePath string `toml:"path"`
	URL      string
	AltText  string `toml:"alt"`
	Width    int
}

type ImgurResponse struct {
	Data ImgurResponseData `json:"data"`
}

type ImgurResponseData struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}
