package completion

import (
	"context"
	"strings"

	hooklist "github.com/kyoh86/gogh/v4/app/hook/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func Hooks(ctx context.Context, svc *service.ServiceSet, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	completions := make([]cobra.Completion, 0)
	for h, err := range hooklist.NewUsecase(svc.HookService).Execute(ctx) {
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		if h == nil {
			continue
		}
		id := h.ID()
		name := h.Name()
		if toComplete != "" &&
			!strings.HasPrefix(id, toComplete) &&
			!strings.HasPrefix(name, toComplete) {
			continue
		}
		completions = append(completions, id+"\t"+name)
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
