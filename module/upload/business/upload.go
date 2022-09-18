package uploadbusiness

import (
	"bytes"
	"context"
	"fmt"
	"food-delivery/common"
	"food-delivery/component/uploadprovider"
	uploadmodel "food-delivery/module/upload/model"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CreateImageStorage interface {
	CreateImage(context context.Context, data *common.Image)
}

type uploadBiz struct {
	provider uploadprovider.UploadProvider
	imgStore CreateImageStorage
}

func NewUploadBiz(provider uploadprovider.UploadProvider, imgStore CreateImageStorage) *uploadBiz {
	return &uploadBiz{
		provider: provider,
		imgStore: imgStore,
	}
}

func (biz *uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string) (*common.Image, error) {

	fileBytes := bytes.NewBuffer(data)

	w, h, err := getImageDimension(fileBytes)

	if err != nil {
		return nil, uploadmodel.ErrFileIsNotImage(err)
	}

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	fileExt := filepath.Ext(fileName)                                // "img.jpg" => ".jpg"
	fileName = fmt.Sprintf("%d%s", time.Now().Nanosecond(), fileExt) // 9129324893248.jpg

	img, err := biz.provider.SaveFileUploaded(ctx, data, fmt.Sprintf("%s/%s", folder, fileName))

	if err != nil {
		return nil, uploadmodel.ErrCannotSaveFile(err)
	}

	img.Width = w
	img.Height = h
	//img.CloudName = "s3" // should be set in provider
	img.Extension = fileExt

	//if err := biz.imgStore.CreateImage(ctx, img); err != nil {
	//	// delete img on S3
	//	return nil, uploadmodel.ErrCannotSaveFile(err)
	//}

	return img, nil
}
func (biz *uploadBiz) UploadLocal(ctx context.Context, data []byte, folder, fileName string) (*common.Image, error) {

	//fileBytes := bytes.NewBuffer(data)

	//w, h, err := getImageDimension(fileBytes)
	//
	//if err != nil {
	//	return nil, uploadmodel.ErrFileIsNotImage(err)
	//}

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	fileExt := filepath.Ext(fileName)                                // "img.jpg" => ".jpg"
	fileName = fmt.Sprintf("%d%s", time.Now().Nanosecond(), fileExt) // 9129324893248.jpg

	img, err := biz.provider.SaveFileUploaded(ctx, data, fmt.Sprintf("%s/%s", folder, fileName))

	if err != nil {
		return nil, uploadmodel.ErrCannotSaveFile(err)
	}

	img.Width = 100
	img.Height = 100
	img.CloudName = "local" // should be set in provider
	img.Extension = fileExt

	//if err := biz.imgStore.CreateImage(ctx, img); err != nil {
	//	// delete img on S3
	//	return nil, uploadmodel.ErrCannotSaveFile(err)
	//}

	return img, nil
}

func getImageDimension(reader io.Reader) (int, int, error) {

	//img, _, err := image.DecodeConfig(reader)
	//fmt.Println(img)
	//if err != nil {
	//	log.Println("err: ", err)
	//	return 0, 0, err
	//}

	return 100, 100, nil
}
func getImageDimension1(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
