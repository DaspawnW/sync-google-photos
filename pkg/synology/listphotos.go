package synology

import (
	"encoding/json"
	"fmt"
)

func (c *Connection) ListPhotosInFolder(folderId int) ([]Photo, error) {
	photoList := []Photo{}
	page := 0

	for {
		l, err := c.listPhotosInFolderPaged(folderId, page)
		if err != nil {
			return nil, err
		}

		photoList = append(photoList, l...)
		if len(l) < PAGE_LIMIT {
			break
		}

		page++
	}

	return photoList, nil
}

func (c *Connection) listPhotosInFolderPaged(folderId int, page int) ([]Photo, error) {
	offset := page * PAGE_LIMIT
	url := fmt.Sprintf("%s/webapi/entry.cgi?api=SYNO.Foto.Browse.Item&version=1&method=list&offset=%d&limit=%d&folder_id=%d&additional=[\"tag\"]&_sid=%s", c.url, offset, PAGE_LIMIT, folderId, c.sid)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	parsedResp := &ListPhotosResponse{}
	derr := json.NewDecoder(resp.Body).Decode(parsedResp)
	if derr != nil {
		return nil, derr
	}

	if parsedResp.Success != true {
		return nil, fmt.Errorf("failed to list photos in folder %d with non successful response", folderId)
	}

	return parsedResp.Data.List, nil
}
