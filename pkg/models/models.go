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
}

type SourceFile struct {
	Metadata SourceFileMetadata `toml:"metadata"`
	Articles []Article          `toml:"articles"`
}

type SourceFileMetadata struct {
	Title       string `toml:"title"`
	IssueNo     int    `toml:"no"`
	PubDate     string `toml:"date"`
	WelcomeText string `toml:"welcome"`
}
