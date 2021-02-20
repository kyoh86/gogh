package command

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
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
				log.FromContext(ctx).WithFields(log.Fields{
					"error":  err,
					"format": project.FullFilePath(),
				}).Info("failed to format")
			}
			fmt.Println(str)
		}
	}
	return nil
}
