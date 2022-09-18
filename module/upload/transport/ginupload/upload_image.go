package ginupload

import (
	"fmt"
	"food-delivery/common"
	"food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
	"image"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func UploadImage(appCtx appctx.AppContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		folder := c.DefaultPostForm("folder", "img")

		fileExt := filepath.Ext(fileHeader.Filename)                      // "img.jpg" => ".jpg"
		fileName := fmt.Sprintf("%d%s", time.Now().Nanosecond(), fileExt) // 9129324893248.jpg

		if err := c.SaveUploadedFile(fileHeader, fmt.Sprintf("static/%s/%s", folder, fileName)); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(common.Image{
			Id:        0,
			Url:       fmt.Sprintf("http://localhost:8080/static/%s/%s", folder, fileName),
			Width:     0,
			Height:    0,
			CloudName: "local",
			Extension: fileExt,
		}))
		//if err != nil {
		//	panic(common.ErrInvalidRequest(err))
		//}
		//
		//folder := c.DefaultPostForm("folder", "img")
		//
		////fmt.Print("folder", folder)
		//
		//file, err := fileHeader.Open()
		//
		//if err != nil {
		//	panic(common.ErrInvalidRequest(err))
		//}
		//
		//defer file.Close() // we can close here
		//
		//dataBytes := make([]byte, fileHeader.Size)
		//if _, err := file.Read(dataBytes); err != nil {
		//	panic(common.ErrInvalidRequest(err))
		//}
		//
		////imgStore := uploadstorage.NewSQLStore(db)
		//biz := uploadbusiness.NewUploadBiz(appCtx.UploadProvider(), nil)
		//img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)
		//
		//if err != nil {
		//	panic(err)
		//}
		//c.JSON(200, common.SimpleSuccessResponse(img))
	}
}

func getImageDimension(reader io.Reader) (int, int, error) {

	img, _, err := image.DecodeConfig(reader)
	fmt.Println(img)
	if err != nil {
		log.Println("err: ", err)
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}
