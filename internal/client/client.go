package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	cfg "github.com/Tihmmm/mr-decorator/internal/config"
	"github.com/Tihmmm/mr-decorator/internal/errors"
	"github.com/doyensec/safeurl"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Client interface {
	GetArtifact(projectId int, jobId int, artifactFileName string, glToken string) (artifactDir string, err error)
	SendNote(note string, projectId int, mergeRequestIid int, glToken string) (err error)
}

type HttpClient struct {
	cfg    cfg.HttpClientConfig
	client *safeurl.WrappedClient
}

func NewHttpClient(cfg cfg.HttpClientConfig) Client {
	config := safeurl.GetConfigBuilder().SetAllowedIPs(cfg.Ip).
		Build()
	httpClient := &HttpClient{
		cfg:    cfg,
		client: safeurl.Client(config),
	}
	return httpClient
}

const (
	jobArtifactsEndpointBasePath      = "/api/v4/projects/%d/jobs/%d/artifacts/%s"
	mergeRequestNotesEndpointBasePath = "/api/v4/projects/%d/merge_requests/%d/notes"
	artifactsBaseDir                  = "artifacts"
	privateTokenHeader                = "PRIVATE-TOKEN"
	contentTypeHeader                 = "Content-Type"
	contentTypeJson                   = "application/json"
)

func (c *HttpClient) GetArtifact(projectId int, jobId int, artifactFileName string, glToken string) (artifactDir string, err error) {
	jobArtifactPath := fmt.Sprintf(jobArtifactsEndpointBasePath, projectId, jobId, artifactFileName)
	req, err := newBaseGetRequest(jobArtifactPath, glToken, c.cfg.Host)
	if err != nil {
		log.Printf("Error creating GET request for job artifact '%s': %s\n", jobArtifactPath, err)
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Error downloading artifact for project: %d, job: %d, err: %s\n", projectId, jobId, err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error downloading artifact for project: %d, job: %d. Gitlab response status: %d\n", projectId, jobId, resp.StatusCode)
		return "", &errors.DownloadError{}
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Printf("Error closing response body for project: %d, job: %d, err: %s\n", projectId, jobId, err)
			return
		}
	}(resp.Body)

	dirPath := filepath.Join(artifactsBaseDir, uuid.New().String())
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		log.Printf("Error creating artifact directory: %s\n", err)
		return "", err
	}
	filePath := filepath.Join(dirPath, artifactFileName)
	out, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating artifact file: %s\n", err)
		return "", err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		log.Printf("Error copying artifact for project: %d, job: %d, err: %v\n", projectId, jobId, err)
		return "", err
	}

	return dirPath, nil
}

func (c *HttpClient) SendNote(note string, projectId int, mergeRequestIid int, glToken string) (err error) {
	body := struct {
		Body string `json:"body"`
	}{note}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error parsing body: %s\n", err)
		return err
	}

	notePath := fmt.Sprintf(mergeRequestNotesEndpointBasePath, projectId, mergeRequestIid)

	req, err := newBasePostRequest(notePath, bytes.NewBuffer(bodyBytes), glToken, c.cfg.Host)
	if err != nil {
		log.Printf("Error creating POST request to send node for job artifact '%s': %v\n", notePath, err)
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Error sending note: %s\n", err)
		return err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Printf("Error closing response body for job artifact '%s': %v\n", notePath, err)
			return
		}
	}(resp.Body)

	var respBuf []byte
	_, err = resp.Body.Read(respBuf)

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Error sending note. Gitlab response status code: %d\nbody: %s\n", resp.StatusCode, string(respBuf))
		return
	}
	return nil
}

func newBaseGetRequest(path string, glToken string, host string) (*http.Request, error) {
	fullLink := "https://" + host + path
	req, err := http.NewRequest(http.MethodGet, fullLink, nil)
	if err != nil {
		log.Printf("Error creating GET request for job artifact '%s': %s\n", path, err)
		return nil, err
	}
	req.Header.Set(privateTokenHeader, glToken)

	return req, nil
}

func newBasePostRequest(path string, body io.Reader, glToken string, host string) (*http.Request, error) {
	fullLink := "https://" + host + path
	req, err := http.NewRequest(http.MethodPost, fullLink, body)
	if err != nil {
		log.Printf("Error creating POST request for job artifact '%s': %s\n", path, err)
		return nil, err
	}
	req.Header.Set(privateTokenHeader, glToken)
	req.Header.Set(contentTypeHeader, contentTypeJson)

	return req, nil
}
