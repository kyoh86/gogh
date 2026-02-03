package completion

import (
	"context"
	"strings"

	overlaylist "github.com/kyoh86/gogh/v4/app/overlay/list"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func Overlays(ctx context.Context, svc *service.ServiceSet, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
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
