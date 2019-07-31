package store

// Article represents a table articles, which needs for simple responses on a specific messages
type Article struct {
	ID   int
	Name string
	Text string
}

// LoadArticle returns message string from actions table
func LoadArticle(name string) (*Article, error) {
	a := &Article{}
	err := dbConn.QueryRow("SELECT text FROM articles WHERE name=$1", name).Scan(&a.Text)
	return a, err
}
