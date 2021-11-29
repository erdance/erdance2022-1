package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/gobuffalo/buffalo-heroku/genny/heroku"
	her "github.com/gobuffalo/buffalo-heroku/heroku"
)

var newOptions = struct {
	*heroku.Options
	dryRun bool
}{
	Options: &heroku.Options{},
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates a new heroku application",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if newOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		opts := newOptions.Options
		var addons her.Addons
		for _, a := range her.DefaultAddons {
			key := "addon:" + a.Name
			v, err := cmd.Flags().GetString(key)
			if err != nil {
				return errors.WithStack(err)
			}
			a.Level = v
			addons = append(addons, a)
		}
		opts.Addons = addons
		gg, err := heroku.New(opts)
		if err != nil {
			return errors.WithStack(err)
		}
		gg.With(r)

		g, err := gogen.Fmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		r.WithRun(func(r *genny.Runner) error {
			cmd := exec.Command("heroku", "open")
			return r.Exec(cmd)
		})

		return r.Run()
	},
}

func init() {
	newCmd.Flags().StringVarP(&newOptions.DynoLevel, "dyno-level", "l", "free", strings.Join(her.DynoLevels, ", "))
	newCmd.Flags().StringVarP(&newOptions.AppName, "app-name", "a", "", "the name of the heroku app to deploy")
	newCmd.Flags().StringVarP(&newOptions.Region, "region", "r", "us", "heroku region to deploy heroku app")
	newCmd.Flags().StringVarP(&newOptions.Environment, "environment", "e", "production", "the environment to run the application in")
	newCmd.Flags().BoolVar(&newOptions.Auth, "auth", false, "log into heroku from the cli")
	newCmd.Flags().BoolVarP(&newOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	for _, a := range her.DefaultAddons {
		newCmd.Flags().String("addon:"+a.Name, a.Level, fmt.Sprintf("options: %s", strings.Join(a.Available, ", ")))
	}
	herokuCmd.AddCommand(newCmd)
}
