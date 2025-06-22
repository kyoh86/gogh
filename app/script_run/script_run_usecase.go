package script_run

import (
	"context"
	"fmt"

	lua "github.com/Shopify/go-lua"
)

// UseCase for running script scripts
type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

// Globals represents a map of global variables to be passed to Lua
type Globals map[string]any

// ToLuaTable converts the Globals map to a Lua table on the given state
func (g Globals) ToLuaTable(l *lua.State) {
	l.NewTable()
	for key, value := range g {
		l.PushString(key)
		switch v := value.(type) {
		case string:
			l.PushString(v)
		case int:
			l.PushInteger(v)
		case float64:
			l.PushNumber(v)
		case bool:
			l.PushBoolean(v)
		case map[string]any:
			// Recursively convert nested maps
			Globals(v).ToLuaTable(l)
		default:
			// For complex types, convert to a simple representation
			l.PushString(fmt.Sprintf("%v", v))
		}
		l.SetTable(-3)
	}
}

type Script struct {
	Code    string
	Globals Globals
}

func (uc *UseCase) Execute(ctx context.Context, script Script) error {
	l := lua.NewState()
	lua.OpenLibraries(l)

	// Set up the global 'gogh' table
	l.NewTable()
	script.Globals.ToLuaTable(l)
	l.SetGlobal("gogh")

	if err := lua.DoString(l, script.Code); err != nil {
		return fmt.Errorf("run Lua: %w", err)
	}
	return nil
}
