package models

type Article struct {
	URL      string
	Title    string
	Type     string
	Category string
	Summary  string
}

type NewsletterData struct {
	Title       string
	IssueNo     int
	PubDate     string
	WelcomeText string
	Articles    []Article
}

type SourceFile struct {
	Metadata SourceFileMetadata  `toml:"metadata"`
	Articles []SourceFileArticle `toml:"articles"`
}

type SourceFileMetadata struct {
	Title       string `toml:"title"`
	IssueNo     int    `toml:"no"`
	PubDate     string `toml:"date"`
	WelcomeText string `toml:"welcome"`
}

type SourceFileArticle struct {
	Url      string `toml:"url"`
	Type     string `toml:"type"`
	Category string `toml:"cat"`
	Summary  string `toml:"summary"`
}
