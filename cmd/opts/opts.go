package opts

import (
	"github.com/Tihmmm/mr-decorator-core/client"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/validator"
)

type CmdOpts struct {
	ConfigPath      string
	ParserConfig    *config.ParserConfig
	DecoratorConfig *config.DecoratorConfig
	C               client.Client
	V               validator.Validator
}
