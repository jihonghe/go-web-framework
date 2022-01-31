package summer

type IGroup interface {
	Get(url string, handler ControllerHandler)
	Post(url string, handler ControllerHandler)
	Put(url string, handler ControllerHandler)
	Delete(url string, handler ControllerHandler)
}
