package inject

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

type Provider interface{}
type Provide []Provider

type Dep map[int]string

type providerInfo struct {
	singleton bool

	name string

	// depends of this provider
	deps []string

	// type of injerct value
	typ reflect.Type

	// provider function
	prov Provider

	// for provider function
	pval reflect.Value
	ptyp reflect.Type

	// instance of provider value
	value interface{}

	// reflect value of provider value
	val reflect.Value

	done uint32
}

func newProvider(obj Object, singleton bool) (info *providerInfo) {
	provVal := reflect.ValueOf(obj.Value)
	provs, ok := obj.Value.(Provide)

	if !ok && provVal.Kind() != reflect.Func {
		info = newProviderValue(obj, provVal, singleton)
		return
	}

	if len(provs) == 0 && provVal.Kind() != reflect.Func {
		panic("empty Provider not permit")
	}

	var depMap Dep
	info = &providerInfo{
		singleton: singleton,
	}

	if len(provs) > 0 {
		info.prov = provs[len(provs)-1]
		depMap, _ = provs[0].(Dep)
		info.pval = reflect.ValueOf(info.prov)
	} else {
		info.prov = obj.Value
		info.pval = provVal
	}

	if !info.pval.IsValid() || info.pval.IsNil() {
		panic("provider can not be nil")
	}

	info.ptyp = info.pval.Type()

	if info.ptyp.Kind() != reflect.Func {
		panic(fmt.Sprintf("expected a func end of provider but get `%v`", info.ptyp))
	}

	if info.ptyp.NumOut() > 0 {
		if obj.Type == nil {
			info.typ = indirectType(info.ptyp.Out(0))
		} else {
			info.typ = indirectType(reflect.TypeOf(obj.Type))
		}

		info.setName(obj.Name)
	}

	numIn := info.ptyp.NumIn()
	deps := make([]string, 0, numIn)

	for i := 0; i < numIn; i++ {
		typ := indirectType(info.ptyp.In(i))
		deps = append(deps, createName(typ, depMap[i]))
	}

	info.deps = deps

	return
}

func newProviderValue(obj Object, val reflect.Value, singleton bool) *providerInfo {
	info := &providerInfo{
		value:     obj.Value,
		singleton: singleton,
	}

	if val.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("provider value must be ptr, but get `%v`", val.Kind()))
	}

	if val.IsNil() {
		panic("provider can not be nil")
	}

	info.val = val

	if obj.Type == nil {
		info.typ = indirectValue(info.val).Type()
	} else {
		info.typ = indirectType(reflect.TypeOf(obj.Type))
	}

	info.setName(obj.Name)

	return info
}

func (p *providerInfo) setName(name string) {
	p.name = createName(p.typ, name)
}

func (p *providerInfo) invoke(inj *injector, status invokeStatus) (out []reflect.Value, err error) {
	inj.clck.RLock()
	if v, ok := inj.caches[p.name]; ok {
		inj.clck.RUnlock()
		out = v
		return
	}
	inj.clck.RUnlock()

	defer func() {
		if err == nil {
			atomic.AddUint32(&p.done, 1)
		}
	}()

	// on process
	status.set(p.name)

	if p.value != nil {
		out = []reflect.Value{p.val}
		return
	}

	in := make([]reflect.Value, 0, len(p.deps))
	for _, dep := range p.deps {
		prov := inj.get(dep)
		if prov == nil {
			err = fmt.Errorf("provider not found of dep <%s> by %v", dep, p.ptyp)
			return
		}

		// avoid cycle dependencies
		if atomic.LoadUint32(&prov.done) == 0 && status.has(prov.name) {
			err = fmt.Errorf("provider cycle dependencies of dep <%s> by %v", dep, p.ptyp)
			return
		}

		// deep invoke
		ot, er := prov.invoke(inj, status)
		if er != nil {
			err = er
			return
		}

		if len(ot) > 0 {
			if !ot[0].IsValid() {
				out = ot
				return
			}
		}

		in = append(in, ot[0])
	}

	// invoke provider function
	out = p.pval.Call(in)
	if p.ptyp.NumOut() > 0 {
		inj.clck.Lock()
		inj.caches[p.name] = out
		inj.clck.Unlock()
	}
	return
}

// create unique name of type
func createName(typ reflect.Type, name string) string {
	return typ.PkgPath() + ":" + typ.Name() + ":" + name
}

// reflect indirect of reflect.Value
func indirectValue(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}

// reflect indirect of reflect.Type
func indirectType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}
