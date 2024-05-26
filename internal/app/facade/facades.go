package facade

// Facades is an interface for facades
type Facades struct {
	AuthorFacade *AuthorFacade
	BookFacade   *BookFacade
	AuthFacade   *AuthFacade
	UserFacade   *UserFacade
}
