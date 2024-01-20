package synology

import "time"

type LoginResponse struct {
	Data struct {
		Did string `json:"did"`
		Sid string `json:"sid"`
	} `json:"data"`
	Success bool `json:"success"`
}

type ListFolderResponse struct {
	Data struct {
		List []Folder `json:"list"`
	} `json:"data"`
	Success bool `json:"success"`
}

type Folder struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ListPhotosResponse struct {
	Data struct {
		List []Photo `json:"list"`
	} `json:"data"`
	Success bool `json:"success"`
}

type TagActionResponse struct {
	Success bool `json:"success"`
}

type Photo struct {
	Id         int    `json:"id"`
	FolderId   int    `json:"folder_id"`
	FileName   string `json:"filename"`
	Time       int    `json:"time"`
	Additional struct {
		Tag []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tag"`
		Thumbnail struct {
			XL       string `json:"xl"`
			CacheKey string `json:"cache_key"`
		} `json:"thumbnail"`
	} `json:"additional"`
}

func (p Photo) HasTag(tagId int) bool {
	for _, t := range p.Additional.Tag {
		if t.Id == tagId {
			return true
		}
	}

	return false
}

func (p Photo) GetTime() time.Time {
	return time.Unix(int64(p.Time), 0)
}
