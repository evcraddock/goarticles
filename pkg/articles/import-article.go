package articles

//ImportArticle represents and article that can be imported
type ImportArticle struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	URL         string   `yaml:"url"`
	Banner      string   `yaml:"banner"`
	Images      []string `yaml:"images"`
	PublishDate string   `yaml:"publishDate"`
	Author      string   `yaml:"author"`
	Categories  []string `yaml:"categories"`
	Tags        []string `yaml:"tags"`
	Layout      string   `yaml:"layout"`
	Content     string   `fm:"content" yaml:"-"`
}
