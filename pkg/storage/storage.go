package storage

// Post - публикация.
type Post struct {
	ID         int
	Title      string
	Content    string
	AuthorID   int
	AuthorName string
	CreatedAt  int64
}

type Author struct {
	ID   int
	Name string
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Posts() ([]Post, error)
	AddPost(Post) error
	UpdatePost(Post) error
	DeletePost(Post) error
	Authors() ([]Author, error)
}
