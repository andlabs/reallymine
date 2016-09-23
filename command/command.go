// 22 october 2015
package command

import (
	"fmt"
	"os"
	"io"
	"strings"
	"reflect"
)

type Command struct {
	Name		string
	Args			[]Arg
	Description	[]string
	Do			interface{}
}

var (
	// this is what text/template does
	typeError = reflect.TypeOf((*error)(nil)).Elem()
)

func (c *Command) validate() (issues []string) {
	var ft reflect.Type

	bad := func(format string, args ...interface{}) {
		issues = append(issues, fmt.Sprintf(format, args...))
		// don't stop testing; there might be multiple issues
	}

	if c.Name == "" {
		bad("name must be specified")
	}
	if strings.Index(c.Name, " ") != -1 {
		bad("name cannot contain spaces")
	}
	if c.Description == "" {
		bad("description must be specified")
	}

	if c.Do == nil {
		bad("function must be specified")
	} else {
		ft = reflect.TypeOf(c.Do)
		if ft.Kind() != reflect.Func {
			bad("not a function")
			ft = nil
		} else {
			if ft.IsVariadic() {
				bad("variadic functions not supported")
			}
			if ft.NumOut() != 1 {
				bad("function must return (error)")
			} else {
				rt := ft.Out(0)
				if rt != typeError {
					bad("function doesn't return error")
				}
			}
			if reflect.ValueOf(c.Do).IsNil() {
				bad("nil function value specified")
			}
			if ft.NumIn() != len(c.Args) {
				bad("function does not take right number of arguments")
				// and clear ft so the following checks don't use it
				ft = nil
			}
		}
	}

	if ft != nil {
		for i, arg := range c.Args {
			t := arg.atype()
			if ft.In(i) != t {
				bad("argument %d not of type %s", i, t)
			}
		}
	}

	return issues
}

var ErrWrongArgCount = fmt.Errorf("wrong number of arguments")

func (c *Command) Invoke(args []string) error {
	var of io.WriteCloser

	if len(args) != len(c.Args) {
		return ErrWrongArgCount
	}
	fv := reflect.ValueOf(c.Do)
	fa := make([]reflect.Value, len(args))
	for i, arg := range c.Args {
		out, err := arg.prepare(args[i])
		if err != nil {
			return err
		}
		defer out.deferfunc()
		fa[i] = out.obj
	}
	out := fv.Call(fa)
	ret := out[0].Interface()
	if ret == nil {
		return nil
	}
	return ret.(error)
}

func Validate(commands []*Command) (problems []string) {
	if len(commands) == 0 {
		return []string{"no commands"}
	}
	for _, c := range commands {
		problems = append(problems, c.validate()...)
	}
	return problems
}

func FormatUsage(commands []*Command) string {
	if len(commands) == 0 {
		// this should not happen, but return something reasonable anyway
		return "(no commands)\n"
	}
	out := ""
	for _, c := range commands {
		out += usageL1(c.Name + arglist(c.Args))
		out += formatDescription(c.Description, c.Args)
	}
	return out
}
