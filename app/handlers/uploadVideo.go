package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ClimenteA/vidcastcutter/app/funcs"
	"github.com/ClimenteA/vidcastcutter/app/kvjson"
	"github.com/ClimenteA/vidcastcutter/app/types"
	"github.com/gofiber/fiber/v2"
)

// const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// func getId() string {

// 	src := rand.NewSource(time.Now().UnixNano())
// 	r := rand.New(src)

// 	b := make([]byte, 8)
// 	for i := range b {
// 		b[i] = letterBytes[r.Intn(len(letterBytes))]
// 	}

// 	return string(b)

// }

var failedToUploadResponse = types.WebResponse{
	Status:  string(types.Success),
	Message: "Failed to upload file",
}

func UploadVideo(c *fiber.Ctx) error {

	kv, _ := c.Locals("kv").(kvjson.DB)

	// Get form
	form, err := c.MultipartForm()
	if err != nil {
		log.Fatal(err)
		return c.JSON(failedToUploadResponse)
	}

	// Make file dir
	// NOTE: maybe in the future allow multiple uploads
	// uploadId := getId()
	// fileDir := "./videos/" + uploadId
	fileDir := "./videos"
	errDir := os.MkdirAll(fileDir+"/frames", os.ModePerm)
	if errDir != nil {
		log.Fatal(errDir)
		return c.JSON(failedToUploadResponse)
	}

	// Save File to dir
	files := form.File["file"]
	for _, file := range files {
		fmt.Println("Uploading file: ", file.Filename, file.Size, file.Header["Content-Type"][0])

		if !strings.HasPrefix(file.Header["Content-Type"][0], "video/") {
			log.Fatal("Not a video format")
			return c.JSON(failedToUploadResponse)
		}

		extension := strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
		filePath := fileDir + "/video." + extension
		errSave := c.SaveFile(file, filePath)
		if errSave != nil {
			log.Fatal(err)
			return c.JSON(failedToUploadResponse)
		}
		go funcs.ExtractFPS(filePath, kv)
	}

	// Return success response
	return c.JSON(types.WebResponse{
		Status:  string(types.Success),
		Message: "File uploaded!",
	})
}
