package repository

// Repositories is an interface for repositories
type Repositories struct {
	BookRepository     *BookRepository
	AuthorRepository   *AuthorRepository
	PropertyRepository *PropertyRepository
	UserRepository     *UserRepository
}
