package links

type ApiLink struct {
	ID         string   `yaml:"id"`
	Title      string   `yaml:"title"`
	URL        string   `yaml:"url"`
	Banner     string   `yaml:"banner"`
	Categories []string `yaml:"categories"`
	Tags       []string `yaml:"tags"`
}
