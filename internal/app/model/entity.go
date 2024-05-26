package model

// Entity is an interface for entity constraints
type Entity interface {
	common | business | Property
}

type Property interface {
	TermOfService
}

type common interface {
	User
}

type business interface {
	Book | Author
}
