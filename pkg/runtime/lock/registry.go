package lock

import (
	"fmt"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/info"
)

const (
	ServiceName = "lock"
)

type Registry interface {
	Register(fs ...*Factory)
	Create(name string) (lock.LockStore, error)
}

type lockRegistry struct {
	stores map[string]func() lock.LockStore
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &lockRegistry{
		stores: make(map[string]func() lock.LockStore),
		info:   info,
	}
}

func (r *lockRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *lockRegistry) Create(name string) (lock.LockStore, error) {
	if f, ok := r.stores[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}