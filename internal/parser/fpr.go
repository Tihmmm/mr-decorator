package parser

import (
	"bufio"
	"fmt"
	"github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/pkg/file"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const (
	fpruCritScriptPath = "./FPRU_crit.sh"
	fpruHighScriptPath = "./FPRU_high.sh"
	criticalCountFile  = "critical_count.txt"
	highCountFile      = "high_count.txt"
	criticalCsv        = "critical.csv"
	highCsv            = "high.csv"
)

type fpr struct {
	highCount       int
	criticalCount   int
	highRecords     []fprRecord
	criticalRecords []fprRecord
}

type fprRecord struct {
	category        string
	path            string
	sscVulnInstance string
}

func ParseFprFile(fileDir string, dest *fpr) (err error) {
	if err := extractVulns(fileDir); err != nil {
		log.Printf("Error parsing fpr: %s\n", err)
		return err
	}

	dest.criticalCount, err = extractVulnCount(filepath.Join(fileDir, criticalCountFile))
	if err != nil {
		log.Printf("Error parsing critical count: %s\n", err)
		return err
	}
	dest.highCount, err = extractVulnCount(filepath.Join(fileDir, highCountFile))
	if err != nil {
		log.Printf("Error parsing high count: %s\n", err)
		return err
	}

	criticalRecords, err := extractRecords(filepath.Join(fileDir, criticalCsv))
	if err != nil {
		return err
	}
	dest.criticalRecords = criticalRecords

	highRecords, err := extractRecords(filepath.Join(fileDir, highCsv))
	if err != nil {
		return err
	}
	dest.highRecords = highRecords

	return nil
}

func extractVulns(fileDir string) error {
	if err := exec.Command(fpruCritScriptPath, fileDir).Run(); err != nil {
		return err
	}
	if err := exec.Command(fpruHighScriptPath, fileDir).Run(); err != nil {
		return err
	}
	return nil
}

func extractVulnCount(filePath string) (int, error) {
	vulnCountFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening vulns count file: %s\n", err)
		return -1, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(vulnCountFile)

	var count int
	scanner := bufio.NewScanner(vulnCountFile)
	for scanner.Scan() {
		lineStr := scanner.Text()
		count, _ = strconv.Atoi(lineStr)
	}
	return count, nil
}

func extractRecords(path string) ([]fprRecord, error) {
	records, err := file.ReadCsv(path)
	if err != nil {
		log.Printf("Error extracting fpr records")
		return []fprRecord{}, err
	}
	var fprRecords []fprRecord
	for i := 1; i < len(records); i++ {
		fprRec := fprRecord{
			category:        records[i][1],
			path:            records[i][2],
			sscVulnInstance: records[i][0],
		}
		fprRecords = append(fprRecords, fprRec)
	}

	return fprRecords, nil
}

func (f *fpr) ToGenSast(cfg config.SastParserConfig, vulnMgmtId int) GenSast {
	var genSast GenSast
	genSast.HcCount = f.vulnCount()
	genSast.HighCount = f.highCount
	genSast.CriticalCount = f.criticalCount
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	genSast.VulnMgmtProjectUrl = baseUrl
	for _, v := range f.highRecords {
		highVulns := Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		genSast.HighVulns = append(genSast.HighVulns, highVulns)
	}
	for _, v := range f.criticalRecords {
		criticalVulns := Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		genSast.CriticalVulns = append(genSast.CriticalVulns, criticalVulns)
	}
	genSast.VulnMgmtReportPath = baseUrl + cfg.ReportPath

	return genSast
}

func (f *fpr) vulnCount() int {
	return f.criticalCount + f.highCount
}
