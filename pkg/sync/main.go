package sync

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daspawnw/sync-google-photos/pkg/googlephotos"
	"github.com/daspawnw/sync-google-photos/pkg/synology"
)

func NewSync(synClient *synology.Connection, googleHttpClient *http.Client, tagId int) Sync {
	return Sync{
		sClient: synClient,
		gClient: googlephotos.GooglePhotosApi{
			OidcClient: googleHttpClient,
		},
		tagId: tagId,
	}
}

type Sync struct {
	sClient *synology.Connection
	gClient googlephotos.GooglePhotosApi
	tagId   int
}

func (s *Sync) Start() error {
	folder, err := s.sClient.ListFolder()
	if err != nil {
		return err
	}

	for _, f := range folder {
		log.Printf("Start sync for folder %s with id %d", f.Name, f.Id)
		files, err := s.sClient.ListPhotosInFolder(f.Id)
		if err != nil {
			return fmt.Errorf("list photos in folder %s with id %d failed with error %v", f.Name, f.Id, err)
		}

		log.Printf("Found %d files in folder %s with id %d", len(files), f.Name, f.Id)

		for _, p := range files {
			if p.HasTag(s.tagId) {
				continue
			}

			path, err := s.sClient.DownloadPhoto(p.Id)
			if err != nil {
				return fmt.Errorf("failed to download photo %d with name %s in folder %s with id %d with error %v", p.Id, p.FileName, f.Name, f.Id, err)
			}

			log.Printf("saved photo %d from folder %s at temp path %s", p.Id, f.Name, *path)
			defer os.Remove(*path)

			uErr := s.gClient.Upload(*path, p.FileName)
			if uErr != nil {
				return fmt.Errorf("failed to upload photo %d from path %s to google with error %v", p.Id, *path, uErr)
			}

			log.Printf("succesfully uploaded photo %d from folder %s to google", p.Id, f.Name)

			tErr := s.sClient.TagPhoto(p.Id, s.tagId)
			if tErr != nil {
				return fmt.Errorf("failed to tag photo with id %d in folder %s with error %v", p.Id, f.Name, tErr)
			}

			log.Printf("succesfully tagged photo %d from folder %s in synology", p.Id, f.Name)
		}
	}

	return nil
}
