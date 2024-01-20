package googlephotos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const UPLOAD_URL = "https://photoslibrary.googleapis.com/v1/uploads"

type GooglePhotosApi struct {
	OidcClient *http.Client
}

func (g *GooglePhotosApi) Upload(filePath string, fileName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	token, err := g.performContentUpload(file, fileName)
	if err != nil {
		return err
	}

	_, cErr := g.createItemFromToken(*token, fileName)
	if cErr != nil {
		return cErr
	}

	return nil
}

func (g *GooglePhotosApi) performContentUpload(file io.ReadCloser, fileName string) (*string, error) {
	req, err := http.NewRequest("POST", UPLOAD_URL, file)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("Content-Length", strconv.FormatInt(upload.size, 10))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Content-Type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-File-Name", fileName)
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	res, err := g.OidcClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("upload failed with invalid status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	sBody := string(body)
	return &sBody, nil
}

func (g *GooglePhotosApi) createItemFromToken(token string, fileName string) (*NewMediaItemResult, error) {
	body := BatchCreateMediaItemsRequest{
		[]NewMediaItem{{
			SimpleMediaItem: SimpleMediaItem{
				UploadToken: token,
				FileName:    fileName,
			},
		}},
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://photoslibrary.googleapis.com/v1/mediaItems:batchCreate", bytes.NewReader(rawBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := g.OidcClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	parsedResp := &BatchCreateMediaItemsResponse{}
	derr := json.NewDecoder(res.Body).Decode(parsedResp)
	if derr != nil {
		return nil, derr
	}

	return &parsedResp.NewMediaItemResults[0], nil
}
