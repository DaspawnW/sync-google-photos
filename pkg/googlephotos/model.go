package googlephotos

type BatchCreateMediaItemsRequest struct {
	NewMediaItems []NewMediaItem `json:"newMediaItems"`
}

type NewMediaItem struct {
	Description     string          `json:"description"`
	SimpleMediaItem SimpleMediaItem `json:"simpleMediaItem"`
}

type SimpleMediaItem struct {
	UploadToken string `json:"uploadToken"`
	FileName    string `json:"fileName"`
}

type BatchCreateMediaItemsResponse struct {
	NewMediaItemResults []NewMediaItemResult `json:"newMediaItemResults"`
}

type NewMediaItemResult struct {
	UploadToken string    `json:"uploadToken"`
	MediaItem   MediaItem `json:"mediaItem"`
}

type MediaItem struct {
	Id              string          `json:"id"`
	Description     string          `json:"description"`
	ProductUrl      string          `json:"productUrl"`
	BaseUrl         string          `json:"baseUrl"`
	MimeType        string          `json:"mimeType"`
	Filename        string          `json:"filename"`
	MediaMetadata   MediaMetadata   `json:"mediaMetadata"`
	ContributorInfo ContributorInfo `json:"contributorInfo"`
}

type MediaMetadata struct {
	CreationTime string `json:"creationTime"`
	Width        string `json:"width"`
	Height       string `json:"height"`
	Photo        Photo  `json:"photo"`
}

type Photo struct {
	CameraMake      string `json:"cameraMake"`
	CameraModel     string `json:"cameraModel"`
	FocalLength     int    `json:"focalLength"`
	ApertureFNumber int    `json:"apertureFNumber"`
	IsoEquivalent   int    `json:"isoEquivalent"`
	ExposureTime    string `json:"exposureTime"`
}

type ContributorInfo struct {
	ProfilePictureBaseUrl string `json:"profilePictureBaseUrl"`
	DisplayName           string `json:"displayName"`
}
