package summer

// NewInstance 定义了如何创建一个新实例，是所有服务类实例的创建入口
type NewInstance func(...interface{}) (interface{}, error)

type ServiceProvider interface {
	// Register 在服务容器中注册一个服务实例的方法，存在【延迟实例化】机制，通过IsDefer()方法判断是否需要延迟实例化
	Register(Container) NewInstance
	// Boot 在调用服务实例时会被调用，负责做一些准备工作：基础配置、初始化参数等
	Boot(Container) error
	// IsDefer 决定在Register服务实例时是否需要执行实例化操作，如果Register时不执行实例化操作，则在make阶段实例化服务
	IsDefer() bool
	// Params 定义传递给NewInstance的参数
	Params(Container) []interface{}
	// Name 服务提供者的凭证
	Name() string
}
