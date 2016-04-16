package imgur

import (
	"net/http"
	"os"
	"testing"

	"github.com/koffeinsource/go-klogger"
)

func TestImageImgurSimulated(t *testing.T) {
	httpC, server := testHTTPClientJSON("{\"data\":{\"id\":\"ClF8rLe\",\"title\":null,\"description\":null,\"datetime\":1451248840,\"type\":\"image\\/jpeg\",\"animated\":false,\"width\":2448,\"height\":3264,\"size\":1071339,\"views\":176,\"bandwidth\":188555664,\"vote\":null,\"favorite\":false,\"nsfw\":null,\"section\":null,\"account_url\":null,\"account_id\":null,\"in_gallery\":false,\"link\":\"http:\\/\\/i.imgur.com\\/ClF8rLe.jpg\"},\"success\":true,\"status\":200}")
	defer server.Close()

	client := new(Client)
	client.HTTPClient = httpC
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = "testing"

	img, status, err := client.GetImageInfo("ClF8rLe")

	if err != nil {
		t.Errorf("GetImageInfo() failed with error: %v", err)
		t.FailNow()
	}

	if img.Animated != false || img.Bandwidth != 188555664 || img.Datetime != 1451248840 || img.Description != "" || img.Height != 3264 || img.Width != 2448 || img.ID != "ClF8rLe" || img.Link != "http://i.imgur.com/ClF8rLe.jpg" || img.Views != 176 {
		t.Fail()
	}

	if status != 200 {
		t.Fail()
	}
}

func TestImageImgurReal(t *testing.T) {
	key := os.Getenv("IMGURCLIENTID")
	if key == "" {
		t.Skip("IMGURCLIENTID environment variable not set.")
	}

	client := new(Client)
	client.HTTPClient = new(http.Client)
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = key

	img, status, err := client.GetImageInfo("ClF8rLe")

	if err != nil {
		t.Errorf("GetRateLimit() failed with error: %v", err)
		t.FailNow()
	}

	if img.Animated != false || img.Datetime != 1451248840 || img.Description != "" || img.Height != 3264 || img.Width != 2448 || img.ID != "ClF8rLe" || img.Link != "http://i.imgur.com/ClF8rLe.jpg" {
		t.Fail()
	}

	if status != 200 {
		t.Fail()
	}
}

func TestImageImgurFailure(t *testing.T) {
	httpC, server := testHTTPClientJSON("{\"data\": {}, \"success\": false, \"status\": 200 }")
	defer server.Close()

	client := new(Client)
	client.HTTPClient = httpC
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = "testing"

	_, _, err := client.GetImageInfo("dummy")

	if err == nil {
		t.Error("GetImageInfo() should have failed, but didn't")
	}
}

func TestImageImgurNotSuccess(t *testing.T) {
	httpC, server := testHTTPClientJSON("{\"data\": {}, \"success\": false, \"status\": 200 }")
	defer server.Close()

	client := new(Client)
	client.HTTPClient = httpC
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = "testing"

	_, _, err := client.GetImageInfo("dummy")

	if err == nil {
		t.Error("GetImageInfo() should have failed, but didn't")
	}
}

func TestImageImgurJSONError(t *testing.T) {
	httpC, server := testHTTPClientInvalidJSON()
	defer server.Close()

	client := new(Client)
	client.HTTPClient = httpC
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = "testing"

	_, _, err := client.GetImageInfo("dummy")

	if err == nil {
		t.Error("GetImageInfo() should have failed, but didn't")
	}
}

func TestImageImgurServerError(t *testing.T) {
	httpC, server := testHTTPClient500()
	defer server.Close()

	client := new(Client)
	client.HTTPClient = httpC
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = "testing"

	_, _, err := client.GetImageInfo("dummy")

	if err == nil {
		t.Error("GetImageInfo() should have failed, but didn't")
	}
}

func TestImageImgurError(t *testing.T) {
	httpC, server := testHTTPClientJSON("{'data' : {}, 'success' : false, 'status'  : 500}")
	defer server.Close()

	client := new(Client)
	client.HTTPClient = httpC
	client.Log = new(klogger.CLILogger)
	client.ImgurClientID = "testing"

	_, _, err := client.GetImageInfo("dummy")

	if err == nil {
		t.Error("GetImageInfo() should have failed, but didn't")
	}
}