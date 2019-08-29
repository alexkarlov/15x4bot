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

// UpdateArticle updates articles by name
func UpdateArticle(name string, t string) error {
	_, err := dbConn.Exec("UPDATE articles SET text=$1 WHERE name=$2", t, name)
	return err
}
