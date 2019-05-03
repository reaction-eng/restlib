package static

/**
Define an interface that all Calc Repos must follow
*/
type Repo interface {
	/**
	Get the public static
	*/
	GetStaticPublicDocument(path string) (string, error)

	/**
	Get the public static
	*/
	GetStaticPrivateDocument(path string) (string, error)
}
