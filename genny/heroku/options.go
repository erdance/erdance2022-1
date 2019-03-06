package heroku

import "github.com/gobuffalo/buffalo-heroku/heroku"

type Options struct {
	AppName     string
	Region      string
	Environment string
	Auth        bool
	DynoLevel   string
	Addons      heroku.Addons
}
