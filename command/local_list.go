package command

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v2"
	"github.com/wacul/ulog"
)

func LocalList(ctx context.Context, roots []string, query string, format gogh.Format) error {
	for _, root := range roots {
		local := gogh.NewLocalController(root)
		projects, err := local.List(ctx, &gogh.LocalListOption{Query: query})
		if err != nil {
			return err
		}
		for _, project := range projects {
			str, err := format(project)
			if err != nil {
				ulog.Logger(ctx).Infof("failed to format %s: %s", project.FullFilePath(), err)
			}
			fmt.Println(str)
		}
	}
	return nil
}
