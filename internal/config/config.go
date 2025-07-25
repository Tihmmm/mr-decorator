package config

import (
	"github.com/caarlos0/env/v10"
	"log"
)

type Config struct {
	Server     ServerConfig
	HttpClient HttpClientConfig
	Parser     ParserConfig
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"3000"`
}
type HttpClientConfig struct {
	Ip   string `env:"GITLAB_IP,unset"`
	Host string `env:"GITLAB_DOMAIN,unset"`
}

type ParserConfig struct {
	ScaParserConfig  ScaParserConfig
	SastParserConfig SastParserConfig
}

type ScaParserConfig struct {
	VulnMgmtProjectUrlTmpl string `env:"SCA_VULN_MGMT_PROJECT_BASE_URL,unset"`
	VulnInstanceTmpl       string `env:"SCA_VULN_MGMT_INSTANCE_SUBPATH_TEMPLATE,unset"`
	ReportPath             string `env:"SCA_VULN_MGMT_REPORT_SUBPATH_TEMPLATE,unset"`
}

type SastParserConfig struct {
	VulnMgmtProjectUrlTmpl string `env:"SAST_VULN_MGMT_PROJECT_BASE_URL,unset"`          // e.g. https://fortify-ssc.company.com/html/ssc/version/%d
	VulnInstanceTmpl       string `env:"SAST_VULN_MGMT_INSTANCE_SUBPATH_TEMPLATE,unset"` // e.g. audit?q=instance_id%3A
	ReportPath             string `env:"SAST_VULN_MGMT_REPORT_SUBPATH_TEMPLATE,unset"`   // e.g. audit?q=analysis_type%3Asca
}

func NewConfig() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error initializing config: %s\n", err)
	}

	return cfg
}
