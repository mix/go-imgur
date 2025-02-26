package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/mix/go-imgur"
)

func printRate(client *imgur.Client) {
	client.Log.Info().Msg("*** RATE LIMIT ***")
	rl, err := client.GetRateLimit()
	if err != nil {
		client.Log.Error().Err(err).Msg("Error in GetRateLimit:")
		return
	}
	client.Log.Info().Msg(fmt.Sprintf("%v", *rl))
}

func printImage(client *imgur.Client, image *string) {
	client.Log.Info().Msg("*** IMAGE ***")
	img, _, err := client.GetImageInfo(*image)
	if err != nil {
		client.Log.Error().Err(err).Msg("Error in GetImageInfo:")
		return
	}
	client.Log.Info().Msg(fmt.Sprintf("%v\n", img))
}

func printAlbum(client *imgur.Client, album *string) {
	client.Log.Info().Msg("*** ALBUM ***")
	img, _, err := client.GetAlbumInfo(*album)
	if err != nil {
		client.Log.Error().Err(err).Msg("Error in GetAlbumInfo:")
		return
	}
	client.Log.Info().Msg(fmt.Sprintf("%v\n", img))
}

func printGImage(client *imgur.Client, gimage *string) {
	client.Log.Info().Msg("*** GALLERY IMAGE ***")
	img, _, err := client.GetGalleryImageInfo(*gimage)
	if err != nil {
		client.Log.Error().Err(err).Msg("Error in GetGalleryImageInfo:")
		return
	}
	client.Log.Info().Msg(fmt.Sprintf("%v\n", img))
}

func printGAlbum(client *imgur.Client, galbum *string) {
	client.Log.Info().Msg("*** GALLERY ALBUM ***")
	img, _, err := client.GetGalleryAlbumInfo(*galbum)
	if err != nil {
		client.Log.Error().Err(err).Msg("Error in GetGalleryAlbumInfo:")
		return
	}
	client.Log.Info().Msg(fmt.Sprintf("%v\n", img))
}

func printURL(client *imgur.Client, url *string) {
	client.Log.Info().Msg("*** URL ***\n")
	img, _, err := client.GetInfoFromURL(*url)
	if err != nil {
		client.Log.Error().Err(err).Msg("Error in GetInfoFromURL:")
		return
	}
	client.Log.Info().Msg(fmt.Sprintf("Image: %+v\n", img.Image))
	client.Log.Info().Msg(fmt.Sprintf("Album: %+v\n", img.Album))
	client.Log.Info().Msg(fmt.Sprintf("GImage: %+v\n", img.GImage))
	client.Log.Info().Msg(fmt.Sprintf("GAlbum: %+v\n", img.GAlbum))
}

func main() {
	imgurClientID := flag.String("id", "", "Your imgur client id. REQUIRED!")
	url := flag.String("url", "", "Gets information based on the URL passed.")
	upload := flag.String("upload", "", "Filepath to an image that will be uploaded to imgur.")
	image := flag.String("image", "", "The image ID to be queried.")
	album := flag.String("album", "", "The album ID to be queried.")
	gimage := flag.String("gimage", "", "The gallery image ID to be queried.")
	galbum := flag.String("galbum", "", "The gallery album ID to be queried.")
	rate := flag.Bool("rate", false, "Get the current rate limit.")
	flag.Parse()

	// Check if there is anything todo
	if flag.NFlag() >= 3 || *imgurClientID == "" {
		flag.PrintDefaults()
		return
	}

	client, err := imgur.NewClient(new(http.Client), *imgurClientID, "")
	if err != nil {
		fmt.Printf("failed during imgur client creation. %+v\n", err)
		return
	}

	if *upload != "" {
		_, st, err := client.UploadImageFromFile(*upload, "", "test title", "test desc")
		if st != 200 || err != nil {
			fmt.Printf("Status: %v\n", st)
			fmt.Printf("Err: %v\n", err)
		}
	}

	if *rate {
		printRate(client)
	}

	if *image != "" {
		printImage(client, image)
	}

	if *album != "" {
		printAlbum(client, album)
	}

	if *gimage != "" {
		printGImage(client, gimage)
	}

	if *galbum != "" {
		printGAlbum(client, galbum)
	}

	if *url != "" {
		printURL(client, url)
	}
}
