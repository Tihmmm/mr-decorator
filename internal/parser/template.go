package parser

const baseTemplateSca = `
<h2>Software composition analysis summary</h2>
{{ if (gt .Count 0) }}
<h3>Vulnerabilities found: {{ .Count }}</h3>
{{ if (gt .Count 10) }}
<i><u>Only showing first 10</i></u><br /><br />
{{ end }}
{{ range $cve := .Cves }}
<details>
<summary><b>{{ $cve.Id }} in {{ $cve.LibraryName }}</b></summary>
<b>Description:</b> {{ $cve.Description }}<br>
{{ if $cve.Recommendations }}
<b>Recommendations:</b> {{ $cve.Recommendations }}
{{ end }}
{{ if $cve.VulnMgmtInstance }} [DETAILS]({{ $cve.VulnMgmtInstance }}) {{ end }}
</details>
{{ end }}
{{ if .VulnMgmtReportPath }}
<h3>[Full report]({{ .VulnMgmtReportPath }})</h3>
{{ end }}
{{ else }}
<h3>No vulnerable components found</h3>
{{ end }}
`

const baseTemplateSast = `
<h2>Static application security testing summary</h2>
{{ if and (.HcCount) (gt .HcCount 0) }}
<h3> Critical vulnerabilities: {{ .CriticalCount }} </h3>
{{ if (gt .CriticalCount 0) }}
{{ if (gt .CriticalCount 10) }}
<i><u>Only showing first 10</i></u><br /><br />
{{ end }}
{{ range $vuln := .CriticalVulns }}
<b>{{ $vuln.Name }}</b> in {{ $vuln.Location }}
{{ if $vuln.VulnMgmtInstance }} [details]({{ $vuln.VulnMgmtInstance }}) {{ end }} <br />
{{ end }}
{{ else }}
<b>None found</b>
{{ end }}
<h3> High vulnerabilities: {{ .HighCount }} </h3>
{{ if (gt .HighCount 0) }}
{{ if (gt .HighCount 10) }}
<i><u>Only showing first 10</i></u><br /><br />
{{ end }}
{{ range $vuln := .HighVulns }}
<b>{{ $vuln.Name }}</b> in {{ $vuln.Location }}
{{ if $vuln.VulnMgmtInstance }} [details]({{ $vuln.VulnMgmtInstance }}) {{ end }} <br />
{{ end }}
{{ else }}
<b>None found</b>
{{ end }}
{{ if .VulnMgmtReportPath }}
<h3> [Full report]({{ .VulnMgmtReportPath }}) <h3>
{{ end }}
{{ else }}
<h3>No serious vulnerabilities found</h3>
{{ end }}
`
