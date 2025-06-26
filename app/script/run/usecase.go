package run

import (
	"context"
	"fmt"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

// Usecase for running script scripts
type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}

// Globals represents a map of global variables to be passed to Lua
type Globals map[string]any

// ToLuaTable converts the Globals map to a Lua table on the given state
func (g Globals) ToLuaTable(l *lua.LState) *lua.LTable {
	table := l.NewTable()
	for key, value := range g {
		switch v := value.(type) {
		case string:
			table.RawSetString(key, lua.LString(v))
		case int:
			table.RawSetString(key, lua.LNumber(v))
		case float64:
			table.RawSetString(key, lua.LNumber(v))
		case bool:
			table.RawSetString(key, lua.LBool(v))
		case map[string]any:
			// Recursively convert nested maps
			table.RawSetString(key, Globals(v).ToLuaTable(l))
		default:
			// For complex types, convert to a simple representation
			table.RawSetString(key, lua.LString(fmt.Sprintf("%v", v)))
		}
	}
	return table
}

type Script struct {
	Code    string
	Globals Globals
}

func (uc *Usecase) Execute(ctx context.Context, script Script) error {
	l := lua.NewState()
	defer l.Close()

	// Load standard libraries
	libs.Preload(l)

	// Set up the global 'gogh' table
	goghTable := script.Globals.ToLuaTable(l)
	l.SetGlobal("gogh", goghTable)

	if err := l.DoString(script.Code); err != nil {
		return fmt.Errorf("run Lua: %w", err)
	}
	return nil
}
