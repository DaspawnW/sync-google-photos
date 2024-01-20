package synology

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dsoprea/go-exif/v3"
	jis "github.com/dsoprea/go-jpeg-image-structure/v2"
)

func (c *Connection) DownloadPhoto(photoId int) (*string, error) {
	photo, err := c.loadPhotoThumbnailInformation(photoId)
	if err != nil {
		return nil, err
	}

	if photo.Additional.Thumbnail.XL != "ready" {
		return nil, fmt.Errorf("failed to download photo with id %d as xl thumbnail is not ready", photoId)
	}

	reader, err := c.loadThumbnail(*photo)
	if err != nil {
		return nil, err
	}

	path, err := c.saveToFile(*photo, reader)
	if err != nil {
		return nil, err
	}

	return path, nil
}

func (c *Connection) loadThumbnail(photo Photo) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/webapi/entry.cgi?api=SYNO.Foto.Thumbnail&version=1&method=get&mode=download&id=%d&type=unit&size=xl&cache_key=%s&_sid=%s", c.url, photo.Id, photo.Additional.Thumbnail.CacheKey, c.sid)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *Connection) saveToFile(p Photo, content io.ReadCloser) (*string, error) {
	file, err := os.CreateTemp("", "photo")
	if err != nil {
		return nil, err
	}

	_, cpErr := io.Copy(file, content)
	if cpErr != nil {
		return nil, cpErr
	}

	fileName := file.Name()

	intfc, err := jis.NewJpegMediaParser().ParseFile(fileName)
	if err != nil {
		return nil, err
	}

	sl := intfc.(*jis.SegmentList)
	ib, err := sl.ConstructExifBuilder()
	if err != nil {
		return nil, err
	}

	ifdExif, err := exif.GetOrCreateIbFromRootIb(ib, "IFD/Exif")
	if err != nil {
		return nil, err
	}

	tm := time.Unix(int64(p.Time), 0)
	setErr := ifdExif.SetStandardWithName("DateTimeOriginal", tm)
	if setErr != nil {
		return nil, setErr
	}

	setExifErr := sl.SetExif(ib)
	if setExifErr != nil {
		return nil, setExifErr
	}

	f, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	sl.Write(f)

	return &fileName, nil
}

func (c *Connection) loadPhotoThumbnailInformation(photoId int) (*Photo, error) {
	url := fmt.Sprintf("%s/webapi/entry.cgi/SYNO.Foto.Browse.Item", c.url)
	body := fmt.Sprintf("api=SYNO.Foto.Browse.Item&method=get&id=[%d]&additional=[\"thumbnail\"]&version=1&_sid=%s", photoId, c.sid)

	b := strings.NewReader(body)
	r, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	parsedResp := &ListPhotosResponse{}
	derr := json.NewDecoder(resp.Body).Decode(parsedResp)
	if derr != nil {
		return nil, derr
	}

	if parsedResp.Success != true || len(parsedResp.Data.List) != 1 {
		return nil, fmt.Errorf("failed to load thumbnail for photo with id %d with non successful response", photoId)
	}

	return &parsedResp.Data.List[0], nil
}
