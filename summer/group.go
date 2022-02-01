package summer

type Group struct {
	core   *Core
	prefix string
}

// NewGroup create the url group
//  param prefix string '/user'
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		prefix: prefix,
	}
}

// Get register a get router
//  url string sub-url '/list'
func (g *Group) Get(url string, handlers ...ControllerHandler) {
	g.core.Get(g.prefix+url, handlers...)
}

// Post register a post router
//  url string sub-url '/add'
func (g *Group) Post(url string, handlers ...ControllerHandler) {
	g.core.Post(g.prefix+url, handlers...)
}

// Put register a put router
//  url string sub-url '/edit'
func (g *Group) Put(url string, handlers ...ControllerHandler) {
	g.core.Put(g.prefix+url, handlers...)
}

// Delete register a delete router
//  url string sub-url '/delete'
func (g *Group) Delete(url string, handlers ...ControllerHandler) {
	g.core.Delete(g.prefix+url, handlers...)
}
