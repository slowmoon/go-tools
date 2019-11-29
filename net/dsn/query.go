package dsn

import (
	"fmt"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

type tagOpt struct {
	Name string
	Default  string
}

func ParseTag(msg string) tagOpt {
	result := strings.SplitN(msg, ",", 2)
	if len(result) == 2 {
		return tagOpt{Name: result[0], Default: result[1]}
	}
	return tagOpt{Name: result[0]}
}


type BindingTypeError struct {
    Name    string
    Type  reflect.Type
}

type InvalidBindingError struct {
	Type reflect.Type
}

func (i InvalidBindingError)Error()string  {
   return fmt.Sprintf("invalid bind type %s", i.Type)
}

func (b *BindingTypeError) Error()string  {
   return fmt.Sprintf("can not decode %s to go type of %s", b.Name, b.Type)
}

type assignFunc   func( t reflect.Value,  opt tagOpt ) error

func stringAssignFunc(value string) assignFunc {
	return func(t reflect.Value, opt tagOpt) error {
		if t.Kind() != reflect.String  || !t.CanSet(){
			return &BindingTypeError{Name:"string", Type: t.Type()  }
		}
		if value == "" {
			t.SetString(opt.Default)
		}  else {
			t.SetString(value)
		}
		return nil
	}
}

func addressesAssignFunc(addrs []string) assignFunc  {
	return func(t reflect.Value, opt tagOpt) error {
		if t.Kind() == reflect.String  {
			if len(addrs)	== 0 &&  opt.Default != "" {
				t.SetString(opt.Default)
			}
			if len(addrs) != 0 {
				t.SetString(addrs[0])
			}
			return  nil
		}
		if !(t.Kind() != reflect.Array || t.Type().Elem().Kind() != reflect.String ) {
			return  &BindingTypeError{Name:"array", Type: t.Type()}
		}
		vals := reflect.MakeSlice(t.Type(), len(addrs), len(addrs))
		for i, addr := range addrs {
			vals.Index(i).SetString(addr)
		}
		if t.CanSet() {
			t.Set(vals)
		}
		return nil
	}
}


type decodeState struct {
    assignFunc   map[string]assignFunc
    used         map[string]bool
    values       url.Values
}

func (d *decodeState)unused() url.Values {
   values :=  url.Values{}
   for k, v := range  d.values {
	   if !d.used[k] {
	   	  values[k] = v
	   }
   }
   return  values
}

func (d *decodeState)decode(v interface{}) (err error)  {
	defer func() {
		if r := recover(); r != nil {
			if  err , ok := r.(runtime.Error);ok {
				panic(err)
			}
			err = r.(error)
		}
	}()

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return &BindingTypeError{Name:"", Type: val.Type()}
	}
	return d.root(val)
}

func (d *decodeState)root(v reflect.Value) (err error) {
    return
}
