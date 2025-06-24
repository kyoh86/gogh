package commands

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/script/run"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptRunCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:    "run",
		Short:  "Run a script gob from stdin (it is internal command used by gogh script-invoke command)",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			var script run.Script
			gob.Register(map[string]any{})
			dec := gob.NewDecoder(os.Stdin)
			if err := dec.Decode(&script); err != nil {
				return fmt.Errorf("decoding script from stdin: %w", err)
			}
			return run.NewUseCase().Execute(ctx, script)
		},
	}
	return cmd, nil
}
