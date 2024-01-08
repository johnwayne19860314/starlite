package inject

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	INJECT_MAX_RECURSIVE_LEVEL = 3
)

type typeError struct {
	err error
}

func (p *typeError) Error() string {
	return p.err.Error()
}

type Injector interface {
	TypeProvider

	// find target object
	Find(interface{}, ...string) error
	Exists(interface{}, ...string) error
	// call provider/func with injected object
	Invoke(interface{}) ([]reflect.Value, error)
	// fill object `inject` tag with providers
	Apply(interface{}) error

	SetParent(Injector) Injector

	// private use
	provInvoker
}

type Object struct {
	Value interface{}
	Type  interface{}
	Name  string
}

type TypeProvider interface {
	SingleProvide(provs ...interface{}) TypeProvider
	SingleProvideAs(prov interface{}, typ interface{}, names ...string) TypeProvider

	Provide(provs ...interface{}) TypeProvider
	ProvideAs(prov interface{}, typ interface{}, names ...string) TypeProvider
}

type provInvoker interface {
	get(string) *providerInfo
}

type invokeStatus map[string]bool

func (m invokeStatus) has(name string) bool {
	_, ok := m[name]
	return ok
}

func (m invokeStatus) set(name string) {
	m[name] = false
}

type injector struct {
	// use for store provider
	values map[string]*providerInfo
	vlck   *sync.RWMutex

	// user for cache provider instance for current inject
	caches map[string][]reflect.Value
	clck   *sync.RWMutex

	parent Injector
}

func New() Injector {
	return &injector{
		values: make(map[string]*providerInfo),
		vlck:   &sync.RWMutex{},
		caches: make(map[string][]reflect.Value),
		clck:   &sync.RWMutex{},
	}
}

func (inj *injector) SingleProvide(provs ...interface{}) TypeProvider {
	for _, prov := range provs {
		inj.SingleProvideAs(prov, nil)
	}
	return inj
}

func (inj *injector) SingleProvideAs(prov interface{}, typ interface{}, names ...string) TypeProvider {
	return inj.provide(prov, typ, true, names...)
}

func (inj *injector) Provide(provs ...interface{}) TypeProvider {
	for _, prov := range provs {
		inj.ProvideAs(prov, nil)
	}
	return inj
}

func (inj *injector) ProvideAs(prov interface{}, typ interface{}, names ...string) TypeProvider {
	return inj.provide(prov, typ, false, names...)
}

func (inj *injector) provide(prov interface{}, typ interface{}, singleton bool, names ...string) TypeProvider {
	var obj Object
	switch p := prov.(type) {
	case Object:
		obj = p
	case *Object:
		obj = *p
	default:
		obj = Object{Value: prov}
	}

	if typ != nil {
		obj.Type = typ
	}
	if len(names) > 0 {
		obj.Name = names[0]
	}

	info := newProvider(obj, singleton)

	// remove exists cache of provider
	inj.clck.Lock()
	delete(inj.caches, info.name)
	inj.clck.Unlock()

	// replace with new prvoder info
	inj.vlck.Lock()
	inj.values[info.name] = info
	inj.vlck.Unlock()
	return inj
}

func (inj *injector) Find(ptr interface{}, names ...string) error {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		panic("need ptr instance")
	}
	provName := createName(indirectType(val.Type()), getName(names...))

	prov := inj.get(provName)
	if prov == nil {
		return fmt.Errorf("provider not found of type `%s`", provName)
	}

	status := make(invokeStatus)

	var (
		out []reflect.Value
		err error
	)
	if prov.singleton && inj.parent != nil {
		out, err = prov.invoke(inj.parent.(*injector), status)
	} else {
		out, err = prov.invoke(inj, status)
	}
	if err != nil {
		return err
	}

	if len(out) > 0 {
		if !out[0].IsValid() {
			return fmt.Errorf("provider value not valid of type `%s`", provName)
		}

		return assignValue(out[0], val.Elem())
	}

	return nil
}

func (inj *injector) Exists(ptr interface{}, names ...string) error {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr {
		panic("need ptr instance")
	}
	provName := createName(indirectType(val.Type()), getName(names...))

	prov := inj.get(provName)
	if prov == nil {
		return fmt.Errorf("provider not found of type `%s`", provName)
	}

	var out []reflect.Value
	if prov.value != nil {
		out = []reflect.Value{prov.val}
	}

	if len(out) > 0 {
		if !out[0].IsValid() {
			return fmt.Errorf("provider value not valid of type `%s`", provName)
		}

		return assignValue(out[0], val.Elem())
	}

	return nil
}

func (inj *injector) Invoke(prov interface{}) ([]reflect.Value, error) {
	status := make(invokeStatus)

	info := newProvider(Object{Value: prov}, false)
	out, err := info.invoke(inj, status)

	if err != nil {
		return nil, fmt.Errorf("provider invoke err: %v", err)
	}

	// remove exists cache of provider
	inj.clck.Lock()
	delete(inj.caches, info.name)
	inj.clck.Unlock()
	return out, nil
}

func (inj *injector) Apply(ptrStruct interface{}) error {
	// status use for check cycle dependencies in current apply flow
	status := make(invokeStatus)

	level := 0

	return inj.apply(ptrStruct, status, level)
}

func (inj *injector) apply(ptrStruct interface{}, status invokeStatus, level int) error {
	level += 1

	val := reflect.ValueOf(ptrStruct)
	elm := reflect.Indirect(val)

	if elm.Kind() != reflect.Struct {
		return &typeError{fmt.Errorf("expected a <*struct> of %v", val)}
	}

	typ := elm.Type()

	for i := 0; i < elm.NumField(); i++ {
		field := elm.Field(i)
		structField := typ.Field(i)

		if !field.CanSet() {
			continue
		}

		tagVal := structField.Tag.Get("inject")
		if tagVal == "-" {
			continue
		}

		if strings.HasPrefix(string(structField.Tag), "inject") || tagVal != "" {
			tagName, tagOpts := parseTag(tagVal)
			// create name of inject value
			provName := createName(indirectType(field.Type()), tagName)

			prov := inj.get(provName)
			if prov == nil {
				if strings.Contains(tagOpts, "omitempty") {
					continue
				}
				return fmt.Errorf("provider not found for type %s:%v", provName, field)
			}

			var (
				out []reflect.Value
				err error
			)
			if prov.singleton && inj.parent != nil {
				out, err = prov.invoke(inj.parent.(*injector), status)
			} else {
				out, err = prov.invoke(inj, status)
			}
			if err != nil {
				return fmt.Errorf("provider invoke of type %s:%v err: %v", provName, field, err)
			}

			if len(out) > 0 {
				if !out[0].IsValid() {
					return fmt.Errorf("value not found for type %s:%v", provName, field)
				}

				if err := assignValue(out[0], field); err != nil {
					return err
				}
			}

			continue
		}

		if level >= INJECT_MAX_RECURSIVE_LEVEL {
			continue
		}

		if field.CanInterface() {
			if field.Kind() == reflect.Struct {
				// restore to pointer struct
				field = field.Addr()
			}

			// child typeError should skip
			if err := inj.apply(field.Interface(), status, level); err != nil {
				if _, ok := err.(*typeError); !ok {
					return err
				}
			}
		}
	}

	return nil
}

func (inj *injector) get(name string) *providerInfo {
	// get provider in current injector
	inj.vlck.Lock()
	if prov := inj.values[name]; prov != nil {
		inj.vlck.Unlock()
		return prov
	}

	inj.vlck.Unlock()

	// back to parent injector
	if inj.parent != nil {
		return inj.parent.get(name)
	}

	return nil
}

// set parent injector
func (inj *injector) SetParent(parent Injector) Injector {
	inj.parent = parent
	return inj
}

func getName(names ...string) string {
	if len(names) > 0 {
		return names[0]
	}
	return ""
}

func parseTag(tag string) (string, string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return tag, ""
}

func assignValue(out, elm reflect.Value) (err error) {
	if out.Type().AssignableTo(elm.Type()) {
		elm.Set(out)
	} else if out.Kind() == reflect.Ptr && out.Type().Elem().AssignableTo(elm.Type()) {
		if out.IsNil() {
			return
		}
		elm.Set(out.Elem())
	} else {
		err = fmt.Errorf("unsupport value assignable from `%v` to `%v`", out.Type(), elm.Type())
	}
	return
}

func GetErrAfterInvoke(outs []reflect.Value) error {
	if len(outs) == 0 {
		return nil
	}

	e := outs[len(outs)-1]
	if !e.IsValid() {
		return nil
	}
	if err, ok := e.Interface().(error); ok {
		return err
	}

	return nil
}
