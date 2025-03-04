package imgur

import (
	"errors"
	"fmt"
	"strings"
)

// GenericInfo is returned from functions for which the final result type is not known beforehand.
// Only one pointer is != nil
type GenericInfo struct {
	Image  *ImageInfo
	Album  *AlbumInfo
	GImage *GalleryImageInfo
	GAlbum *GalleryAlbumInfo
	Limit  *RateLimit
}

// GetInfoFromURL tries to query imgur based on information identified in the URL.
// returns image/album info, http status code of the request, error
func (client *Client) GetInfoFromURL(url string) (*GenericInfo, int, error) {
	url = strings.TrimSpace(url)

	// https://i.imgur.com/<id>.jpg -> image
	if strings.Contains(url, "://i.imgur.com/") {
		return client.directImageURL(url)
	}

	// https://imgur.com/a/<id> -> album
	if strings.Contains(url, "://imgur.com/a/") || strings.Contains(url, "://m.imgur.com/a/") {
		return client.albumURL(url)
	}

	// https://imgur.com/gallery/<id> -> gallery album
	if strings.Contains(url, "://imgur.com/gallery/") || strings.Contains(url, "://m.imgur.com/gallery/") {
		return client.galleryURL(url)
	}

	// https://imgur.com/<id> -> image
	if strings.Contains(url, "://imgur.com/") || strings.Contains(url, "://m.imgur.com/") {
		return client.imageURL(url)
	}

	return nil, -1, errors.New("URL pattern matching for URL " + url + " failed.")
}

func (client *Client) directImageURL(url string) (*GenericInfo, int, error) {
	var ret GenericInfo
	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, ".")
	if start+1 >= end {
		return nil, -1, errors.New("Could not find ID in URL " + url + ". I was going down i.imgur.com path.")
	}
	id := url[start:end]
	client.Log.Debug().Msg(fmt.Sprintf("Detected imgur image ID %v. Was going down the i.imgur.com/ path.", id))
	gii, status, err := client.GetGalleryImageInfo(id)
	if err == nil && status < 400 {
		ret.GImage = gii
	} else {
		var ii *ImageInfo
		ii, status, err = client.GetImageInfo(id)
		ret.Image = ii
	}
	return &ret, status, err
}

func (client *Client) albumURL(url string) (*GenericInfo, int, error) {
	var ret GenericInfo

	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, "?")
	if end == -1 {
		end = len(url)
	}
	id := url[start:end]
	if id == "" {
		return nil, -1, errors.New("Could not find ID in URL " + url + ". I was going down imgur.com/a/ path.")
	}
	client.Log.Debug().Msg(fmt.Sprintf("Detected imgur album ID %v. Was going down the imgur.com/a/ path.", id))
	ai, status, err := client.GetAlbumInfo(id)
	ret.Album = ai
	return &ret, status, err
}

func (client *Client) galleryURL(url string) (*GenericInfo, int, error) {
	var ret GenericInfo

	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, "?")
	if end == -1 {
		end = len(url)
	}
	id := url[start:end]
	if id == "" {
		return nil, -1, errors.New("Could not find ID in URL " + url + ". I was going down imgur.com/gallery/ path.")
	}
	client.Log.Debug().Msg(fmt.Sprintf("Detected imgur gallery ID %v. Was going down the imgur.com/gallery/ path.", id))
	ai, status, err := client.GetGalleryAlbumInfo(id)
	if err == nil && status < 400 {
		ret.GAlbum = ai
		return &ret, status, nil
	}
	// fallback to GetGalleryImageInfo
	client.Log.Debug().Msg(fmt.Sprintf("Failed to retrieve imgur gallery album. Attempting to retrieve imgur gallery image. err: %v status: %d", err, status))
	ii, statusCode, err := client.GetGalleryImageInfo(id)
	if err != nil {
		return nil, statusCode, fmt.Errorf("client.GetGalleryImageInfo:%w", err)
	}
	ret.GImage = ii
	return &ret, statusCode, nil
}

func (client *Client) imageURL(url string) (*GenericInfo, int, error) {
	var ret GenericInfo

	start := strings.LastIndex(url, "/") + 1
	end := strings.LastIndex(url, "?")
	if end == -1 {
		end = len(url)
	}
	id := url[start:end]
	if id == "" {
		return nil, -1, errors.New("Could not find ID in URL " + url + ". I was going down imgur.com/ path.")
	}
	// check if id is a full filename E.G `vsadghes.jpg`, and if so, extract the actual id `vsadghes`
	hasDotIndex := strings.LastIndex(id, ".")
	if hasDotIndex > -1 {
		id = id[0:hasDotIndex]
	}
	client.Log.Debug().Msg(fmt.Sprintf("Detected imgur image ID %v. Was going down the imgur.com/ path.", id))
	ii, status, err := client.GetGalleryImageInfo(id)
	if err == nil && status < 400 {
		ret.GImage = ii
		return &ret, status, nil
	}

	i, statusCode, err := client.GetImageInfo(id)
	if err != nil {
		return nil, statusCode, fmt.Errorf("client.GetImageInfo:%w", err)
	}
	ret.Image = i
	return &ret, statusCode, nil
}
