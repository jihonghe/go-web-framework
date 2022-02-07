package summer

import (
	"errors"
	"fmt"
	"sync"
)

type Container interface {
	// BindSrvProvider 绑定/替换一个服务提供者，如果服务凭证已存在，则会执行替换操作
	BindSrvProvider(provider ServiceProvider) error
	// IsBindSrvProvider 服务凭证是否已绑定了服务提供者
	IsBindSrvProvider(key string) bool

	// Make 根据服务凭证获取服务提供者，如果服务凭证未绑定对应的服务提供者，则err!=nil
	Make(key string) (interface{}, error)
	// MustMake 根据服务凭证获取服务提供者，如果未绑定则panic
	MustMake(key string) interface{}
	// MakeNew 根据传递的参数，创建一个新的服务提供者
	MakeNew(key string, params ...interface{}) (interface{}, error)
}

type SummerContainer struct {
	Container
	providers    map[string]ServiceProvider
	instances    map[string]interface{}
	sync.RWMutex // 为了控制并发
}

func NewSummerContainer() *SummerContainer {
	return &SummerContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
		RWMutex:   sync.RWMutex{},
	}
}

func (s *SummerContainer) PrintProviders() []string {
	res := make([]string, 0)
	for _, provider := range s.providers {
		res = append(res, fmt.Sprint(provider.Name()))
	}
	return res
}

func (s *SummerContainer) findSrvProvider(key string) ServiceProvider {
	if provider, ok := s.providers[key]; ok {
		return provider
	}

	return nil
}

func (s *SummerContainer) BindSrvProvider(provider ServiceProvider) error {
	s.Lock()
	defer s.Unlock()
	key := provider.Name()
	// 更新/替换
	s.providers[key] = provider

	// 如果服务已经实例化，则可以执行配置，参数等的加载，以及注册行为
	if !provider.IsDefer() {
		if err := provider.Boot(s); err != nil {
			return err
		}
		params := provider.Params(s)
		method := provider.Register(s)
		instance, err := method(params)
		if err != nil {
			return err
		}
		s.instances[key] = instance
	}
	return nil
}

func (s *SummerContainer) IsBindSrvProvider(key string) bool {
	return s.findSrvProvider(key) != nil
}

func (s *SummerContainer) Make(key string) (interface{}, error) {
	return s.make(key, nil, false)
}

func (s *SummerContainer) MustMake(key string) interface{} {
	provider, err := s.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return provider
}

func (s *SummerContainer) MakeNew(key string, params ...interface{}) (interface{}, error) {
	return s.make(key, params, true)
}

// make 创建/返回已有服务实例
func (s *SummerContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	s.Lock()
	defer s.Unlock()

	// 查询key对应的服务提供者是否已注册，若未注册则返回错误
	provider := s.findSrvProvider(key)
	if provider == nil {
		return nil, errors.New("unregistered contract: " + key)
	}

	// 重新从创建新服务实例
	if forceNew {
		return s.newInstance(provider, params)
	}

	// 如果容器中已存在服务实例，则直接返回该实例即可
	if ins, ok := s.instances[key]; ok {
		return ins, nil
	}

	// 容器中key对应的服务尚未实例化(延迟实例化机制导致)，则创建新的服务实例
	ins, err := s.newInstance(provider, params)
	if err != nil {
		return nil, err
	}
	// 将创建的服务实例放到容器中
	s.instances[key] = ins
	return ins, nil
}

func (s *SummerContainer) newInstance(provider ServiceProvider, params []interface{}) (interface{}, error) {
	if err := provider.Boot(s); err != nil {
		return nil, err
	}

	if params == nil {
		params = provider.Params(s)
	}
	method := provider.Register(s)
	ins, err := method(params...)
	if err != nil {
		return nil, err
	}
	return ins, nil
}
