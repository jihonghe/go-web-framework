package summer

type ControllerHandler func(c *Context) error

type handlerPipeline struct {
	handlers []ControllerHandler // handler链路
	index    int                 // 下一个要执行的handler的下标
}

func newHandlerPipeline() *handlerPipeline {
	return &handlerPipeline{}
}

func (h *handlerPipeline) next(ctx *Context) error {
	if h.index >= len(h.handlers) {
		return nil
	}
	err := h.handlers[h.index](ctx)
	h.index++
	return err
}
