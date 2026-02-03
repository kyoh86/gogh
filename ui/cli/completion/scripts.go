package completion

import (
	"context"
	"strings"

	scriptlist "github.com/kyoh86/gogh/v4/app/script/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func Scripts(ctx context.Context, svc *service.ServiceSet, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	completions := make([]cobra.Completion, 0)
	for s, err := range scriptlist.NewUsecase(svc.ScriptService).Execute(ctx) {
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		if s == nil {
			continue
		}
		id := s.ID()
		name := s.Name()
		if toComplete != "" &&
			!strings.HasPrefix(id, toComplete) &&
			!strings.HasPrefix(name, toComplete) {
			continue
		}
		completions = append(completions, cobra.Completion(id+"\t"+name))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
