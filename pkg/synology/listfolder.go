package synology

import (
	"encoding/json"
	"fmt"
)

func (c *Connection) ListFolder() ([]Folder, error) {
	page := 0
	folderList := []Folder{}

	for {
		l, err := c.listPagedFolder(-1, page)
		if err != nil {
			return nil, err
		}

		folderList = append(folderList, l...)
		if len(l) < PAGE_LIMIT {
			break
		}

		page++
	}

	// for i loop to get also new appended elements for a recursive search;
	// range with object doesn't perform the range over the newly added elements
	for i := 0; i < len(folderList); i++ {
		f := folderList[i]

		l, err := c.listChildFolder(f.Id)
		if err != nil {
			return nil, err
		}

		folderList = append(folderList, l...)
	}

	return folderList, nil

}

func (c *Connection) listPagedFolder(id int, page int) ([]Folder, error) {
	offset := page * PAGE_LIMIT
	url := fmt.Sprintf("%s/webapi/entry.cgi?api=SYNO.Foto.Browse.Folder&version=1&method=list&offset=%d&limit=%d&_sid=%s", c.url, offset, PAGE_LIMIT, c.sid)
	if id > 0 {
		url += fmt.Sprintf("&id=%d", id)
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	parsedResp := &ListFolderResponse{}
	derr := json.NewDecoder(resp.Body).Decode(parsedResp)
	if derr != nil {
		return nil, derr
	}

	if parsedResp.Success != true {
		return nil, fmt.Errorf("failed to list folder of images non successful response")
	}

	return parsedResp.Data.List, nil
}

func (c *Connection) listChildFolder(id int) ([]Folder, error) {
	page := 0
	folderList := []Folder{}

	for {
		l, err := c.listPagedFolder(id, page)
		if err != nil {
			return nil, err
		}

		folderList = append(folderList, l...)
		if len(l) < PAGE_LIMIT {
			break
		}

		page++
	}

	return folderList, nil
}
