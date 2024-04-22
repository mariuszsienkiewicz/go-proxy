package redirect

type Cache interface {
	Add(hash string, target string)
	Find(hash string) (string, bool)
	Clear()
}
