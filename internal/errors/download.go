package errors

type DownloadError struct {
}

func (e *DownloadError) Error() string {
	return "can't download artifact"
}
