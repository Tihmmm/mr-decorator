package parser

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator/internal/config"
)

type Sca interface {
	vulnCount() int
	ToGenSca(cfg config.ScaParserConfig, vulnMgmtId int) GenSca
}

type GenSca struct {
	Count              int
	Cves               []Cve
	VulnMgmtProjectUrl string
	VulnMgmtReportPath string
}

type Cve struct {
	Id               string
	LibraryName      string
	Description      string
	Recommendations  string
	VulnMgmtInstance string
}

func (sca *GenSca) ApplyLimit() {
	if sca.Count > scaVulnsDisplayed {
		sca.Cves = sca.Cves[:scaVulnsDisplayed]
	}
}

type dependencyCheck struct {
	Dependencies []struct {
		LibraryName     string `json:"fileName"`
		Vulnerabilities []struct {
			CveId       string `json:"name"`
			Description string `json:"description"`
		} `json:"vulnerabilities,omitempty"`
	} `json:"dependencies"`
}

func (dc *dependencyCheck) ToGenSca(cfg config.ScaParserConfig, vulnMgmtId int) GenSca {
	var genSca GenSca
	genSca.Count = dc.vulnCount()
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	genSca.VulnMgmtProjectUrl = baseUrl
	for _, v := range dc.Dependencies {
		for _, vuln := range v.Vulnerabilities {
			cve := Cve{
				Id:          vuln.CveId,
				LibraryName: v.LibraryName,
				Description: vuln.Description,
			}
			genSca.Cves = append(genSca.Cves, cve)
		}
	}
	genSca.VulnMgmtReportPath = baseUrl + cfg.ReportPath

	return genSca
}

func (dc *dependencyCheck) vulnCount() int {
	var count int
	for i := range dc.Dependencies {
		count += len(dc.Dependencies[i].Vulnerabilities)
	}

	return count
}

type cycloneDX struct {
	Vulnerabilities []struct {
		CveId          string `json:"id"`
		Description    string `json:"description"`
		Recommendation string `json:"recommendation"`
		Affects        []struct {
			LibraryName string `json:"ref"`
		} `json:"affects"`
	} `json:"vulnerabilities"`
}

func (dx *cycloneDX) ToGenSca(cfg config.ScaParserConfig, vulnMgmtId int) GenSca {
	var genSca GenSca
	genSca.Count = dx.vulnCount()
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	genSca.VulnMgmtProjectUrl = baseUrl
	for _, vuln := range dx.Vulnerabilities {
		cve := Cve{
			Id:              vuln.CveId,
			LibraryName:     vuln.Affects[0].LibraryName,
			Description:     vuln.Description,
			Recommendations: vuln.Recommendation,
		}
		genSca.Cves = append(genSca.Cves, cve)
	}
	genSca.VulnMgmtReportPath = baseUrl + cfg.ReportPath

	return genSca
}

func (dx *cycloneDX) vulnCount() int {
	return len(dx.Vulnerabilities)
}
