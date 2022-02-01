package summer

type ControllerHandler func(c *Context) error

type handlerPipeline struct {
	handlers []ControllerHandler // handler链路
	index    int                 // 要执行的控制器的下标
}

func newHandlerPipeline() *handlerPipeline {
	return &handlerPipeline{
		index: -1, // 控制器调用链的起始值为0，故index初始值设为-1
	}
}

func (h *handlerPipeline) add(handlers []ControllerHandler) {
	h.handlers = append(h.handlers, handlers...)
}

func (h *handlerPipeline) next(ctx *Context) error {
	// 为什么要先执行自增？因为这是个控制器调用链，只有先自增才能保证调用链中的所有控制器都被执行到
	h.index++
	if h.index < len(h.handlers) {
		return h.handlers[h.index](ctx)
	}

	// 为什么在此执行自增不行？因为控制器调用链是层层调用的，如果在调用函数之后再执行，则会导致调用链中的其他控制器未被执行到
	// h.index++

	return nil
}
