package commands

import (
	"context"
	"strings"

	hooklist "github.com/kyoh86/gogh/v4/app/hook/list"
	overlaylist "github.com/kyoh86/gogh/v4/app/overlay/list"
	scriptlist "github.com/kyoh86/gogh/v4/app/script/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func completeScripts(ctx context.Context, svc *service.ServiceSet, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
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

func completeOverlays(ctx context.Context, svc *service.ServiceSet, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	completions := make([]cobra.Completion, 0)
	for o, err := range overlaylist.NewUsecase(svc.OverlayService).Execute(ctx) {
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		if o == nil {
			continue
		}
		id := o.ID()
		name := o.Name()
		if toComplete != "" &&
			!strings.HasPrefix(id, toComplete) &&
			!strings.HasPrefix(name, toComplete) {
			continue
		}
		completions = append(completions, cobra.Completion(id+"\t"+name))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeHooks(ctx context.Context, svc *service.ServiceSet, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
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
		completions = append(completions, cobra.Completion(id+"\t"+name))
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}
