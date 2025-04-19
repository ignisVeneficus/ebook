package eBookData

type Metadata interface {
	Author() []string
	Title() string
	Publisher() string
	PubDate() string
	ISBN() string
	Contributor() []string
}

type Book interface {
	Metadata() Metadata
	Cover() []byte
}
