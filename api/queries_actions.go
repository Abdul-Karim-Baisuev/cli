package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cli/cli/v2/internal/ghinstance"
	"github.com/cli/cli/v2/internal/ghrepo"
)

func (c Client) JobLog(repo ghrepo.Interface, jobID string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%srepos/%s/actions/jobs/%s/logs",
		ghinstance.RESTPrefix(repo.RepoHost()), ghrepo.FullName(repo), jobID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("job not found")
	} else if resp.StatusCode != 200 {
		return nil, HandleHTTPError(resp)
	}

	return resp.Body, nil
}
