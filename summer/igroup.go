package summer

type IGroup interface {
	Get(url string, handlers ...ControllerHandler)
	Post(url string, handlers ...ControllerHandler)
	Put(url string, handlers ...ControllerHandler)
	Delete(url string, handlers ...ControllerHandler)
}
