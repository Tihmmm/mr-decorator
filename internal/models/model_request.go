package models

type MRRequest struct {
	AuthToken           string `json:"auth_token" validate:"required"`
	FilePath            string
	ProjectId           int    `json:"project_id" validate:"required"`
	JobId               int    `json:"job_id" validate:"required"`
	ArtifactFormat      string `json:"format" validate:"required"`
	ArtifactFileName    string `json:"artifact_file_name" validate:"required"`
	MergeRequestIid     int    `json:"merge_request_iid" validate:"required"`
	VulnerabilityMgmtId int    `json:"vulnerability_mgmt_id" validate:"required"`
}

const (
	FprFn                 = "current.fpr"
	CyclonedxJsonFn       = "cdx.json"             // "http://cyclonedx.org/schema/bom-1.5.schema.json"
	DependencyCheckJsonFn = "depcheck-report.json" // "reportSchema": "1.1"
)
