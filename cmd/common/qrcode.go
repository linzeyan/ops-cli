package common

import (
	"image"
	"image/color"

	/* avoid image: unknown format. */
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/skip2/go-qrcode"

	"github.com/makiuchi-d/gozxing"
	qrread "github.com/makiuchi-d/gozxing/qrcode"
)

/* Generate QR code by giving content. */
func GenerateQRCode(content string, size int, dest string) error {
	if size <= 10 {
		size = 600
	}
	if dest == "" {
		dest = "./qrcode.png"
	}
	if content == "" {
		return ErrInvalidArg
	}
	return qrcode.WriteColorFile(content, qrcode.Medium, size, color.White, color.Black, dest)
}

func ReadQRCode(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	png, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	qrReader := qrread.NewQRCodeReader()
	result, err := qrReader.Decode(png, nil)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
