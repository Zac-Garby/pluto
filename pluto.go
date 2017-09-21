package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Zac-Garby/pluto/bytecode"
	"github.com/Zac-Garby/pluto/compiler"
	"github.com/Zac-Garby/pluto/object"
	"github.com/Zac-Garby/pluto/parser"
	"github.com/Zac-Garby/pluto/vm"

	"github.com/fatih/color"
)

func main() {
	store := vm.NewStore()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")

		if obj, err := execute(text, "<repl>", store); err != nil {
			color.Red("  %s", err)
		} else if obj != nil {
			color.Cyan("  %s", obj)
		}
	}
}

func execute(text, file string, store *vm.Store) (object.Object, error) {
	var (
		cmp   = compiler.New()
		parse = parser.New(text, file)
		prog  = parse.Parse()
	)

	if len(parse.Errors) > 0 {
		parse.PrintErrors()
		return nil, nil
	}

	err := cmp.CompileProgram(prog)
	if err != nil {
		return nil, err
	}

	code, err := bytecode.Read(cmp.Bytes)
	if err != nil {
		return nil, err
	}

	store.Names = cmp.Names
	store.FunctionStore.Define(cmp.Functions...)
	store.Patterns = cmp.Patterns

	fmt.Println(store.Names, cmp.Constants)

	machine := vm.New()
	machine.Run(code, store, cmp.Constants, true)

	if machine.Error != nil {
		return nil, machine.Error
	}

	return machine.ExtractValue(), nil
}
