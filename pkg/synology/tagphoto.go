package synology

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (c *Connection) TagPhoto(photoId int, tagId int) error {
	url := fmt.Sprintf("%s/webapi/entry.cgi/SYNO.Foto.Browse.Item", c.url)
	body := fmt.Sprintf("api=SYNO.Foto.Browse.Item&method=add_tag&version=1&id=[%d]&tag=[%d]&_sid=%s", photoId, tagId, c.sid)

	b := strings.NewReader(body)
	r, err := http.NewRequest("POST", url, b)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	parsedResp := &TagActionResponse{}
	derr := json.NewDecoder(resp.Body).Decode(parsedResp)
	if derr != nil {
		return derr
	}

	if parsedResp.Success == true {
		return nil
	}

	return fmt.Errorf("failed to tag photo with id %d with non successful response from synology nas", photoId)
}
