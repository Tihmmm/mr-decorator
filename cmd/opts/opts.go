package opts

import (
	"github.com/Tihmmm/mr-decorator-core/client"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/validator"
)

type CmdOpts struct {
	Cfg *config.Config
	C   client.Client
	V   validator.Validator
}
