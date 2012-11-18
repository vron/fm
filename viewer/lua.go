package main

import (
	"bufio"
	//"fmt"
	lua "github.com/xenith-studios/golua"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
)

var sesLog *log.Logger
var shell *log.Logger
var luastate *lua.State

func initLua() {
	// Create a new log file
	f, e := os.Create("session.lua")
	if e != nil {
		log.Fatal(e)
	}
	sesLog = log.New(f, "", 0)
	shell = log.New(os.Stdout, "", 0)

	sesLog.Println("-- Session started ", time.Now())

	L := lua.NewState()
	L.OpenLibs()
	luastate = L

	// Register all our functions
	register(L)

	// We should now capture all input
	go readInput(L)
}

func runFile(p string) {
	f, e := ioutil.ReadFile(p)
	if e != nil {
		log.Println("Could not read file: ", e)
		return
	}
	runString(string(f))
}

func runString(s string) {
	l := luastate
	er := l.LoadString(s)
	if er != 0 {
		// Pop of the error description
		str := l.ToString(-1)
		shell.Println("Error loading: ", str)
		return
	}
	er = l.PCall(0, lua.LUA_MULTRET, 0)
	if er != 0 {
		str := l.ToString(-1)
		shell.Println("Error running: ", str)
		return
	}
	sesLog.Println(s)
}

func readInput(l *lua.State) {
	bstdin := bufio.NewReader(os.Stdin)
	for {
		var line string
		// TODO: Something better for windows?
		line, e := bstdin.ReadString('\n')
		if e != nil {
			shell.Println(e)
		}
		er := l.LoadString(line)
		if er != 0 {
			// Pop of the error description
			str := l.ToString(-1)
			shell.Println("Error loading: ", str)
			continue
		}
		er = l.PCall(0, lua.LUA_MULTRET, 0)
		if er != 0 {
			str := l.ToString(-1)
			shell.Println("Error running: ", str)
			continue
		}
		sesLog.Println(line)
	}
}

func register(l *lua.State) {
	for k, v := range funcs {
		// k is ideftifier string and v is the function to call

		// First we check the function
		t := reflect.TypeOf(v)
		if t.Kind() != reflect.Func {
			log.Println("Invalid function: ", k)
			continue
		}

		f, e := createwrap(reflect.ValueOf(v))

		if e != nil {
			log.Println("Error creating wrapper '", k, "': ", e)
			continue
		}

		l.Register(k, f)
	}
}

// Should only be called from register!
func createwrap(f reflect.Value) (func(*lua.State) int, error) {
	/*
		f := reflect.ValueOf(v)
		if len(params) != f.Type().NumIn() {
			err = errors.New("The number of params is not adapted.")
			return
		}
		in := make([]reflect.Value, len(params))
		for k, param := range params {
			in[k] = reflect.ValueOf(param)
		}
		result = f.Call(in)
	*/
	return func(l *lua.State) int {
		// Go through all the in parameters and pop them from the stack
		noargs := l.GetTop()
		if noargs != f.Type().NumIn() {
			log.Println("Non matching number of arguments")
			return 0
		}

		// Build input
		in := make([]reflect.Value, noargs)
		for i := 0; i < noargs; i++ {
			ka := f.Type().In(i)
			switch ka.Kind() {
			case reflect.Uint32:
				if !l.IsNumber(i + 1) {
					log.Println("Expect numeric type ", i)
					return 0
				}
				in[i] = reflect.ValueOf(uint32(l.ToNumber(i + 1)))
			case reflect.Float64:
				if !l.IsNumber(i + 1) {
					log.Println("Expect numeric type ", i)
					return 0
				}
				in[i] = reflect.ValueOf(l.ToNumber(i + 1))
			default:
				log.Println("Unknown argument type ", i)
				return 0
			}
		}

		result := f.Call(in)

		// And return the return values to lua
		for i := 0; i < len(result)-1; i++ {
			v := result[i]
			switch v.Kind() {
			case reflect.Uint32:
				l.PushNumber(float64(v.Uint()))
			default:
				log.Println("Unknown return type ", i)
				return 0
			}
		}

		// We expect the last result to be a error and it to be nill..
		if len(result) < 1 {
			log.Println("Must at least haave on return value, error")
			return 0
		}
		lr := result[len(result)-1]
		if lr.Kind() != reflect.Interface {
			log.Println("Last return should be a error value")
			return 0
		}
		var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
		if !lr.Type().Implements(typeOfError) {
			log.Println("Does not implement error...")
			return 0
		}
		e, ok := lr.Interface().(error)
		if ok {
			log.Println("We have a error: ", e)
			return 0
		}

		return len(result) - 1
	}, nil
}

// "test": func(*lua.State) int { return 0 },

// Please don't store anything but functions here...
var funcs = map[string]interface{}{
	"node": Node,
	"quad": Quad,
	"tria": Tria,
}
