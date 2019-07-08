package gen

type File struct {
	Path   string
	Source Interface
}

type Bundle struct {
	Files []File
}

func NewBundle() *Bundle {
	return &Bundle{}
}

func (b *Bundle) Add(path string, src Interface) {
	b.Files = append(b.Files, File{
		Path:   path,
		Source: src,
	})
}
