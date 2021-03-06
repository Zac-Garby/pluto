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
	"github.com/Zac-Garby/pluto/store"
	"github.com/Zac-Garby/pluto/vm"

	"github.com/fatih/color"
)

func main() {
	store := store.New()
	usePrelude := true

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")

		if obj, err := execute(text, "<repl>", store, usePrelude); err != nil {
			color.Red("  %s", err)
		} else if obj != nil {
			color.Cyan("  %s", obj)
		}

		usePrelude = false
	}
}

func execute(text, file string, store *store.Store, prelude bool) (object.Object, error) {
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

	machine := vm.New()
	machine.Run(code, store, cmp.Constants, prelude)

	if machine.Error != nil {
		return nil, machine.Error
	}

	return machine.ExtractValue(), nil
}
