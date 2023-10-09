package utils

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/sirupsen/logrus"
	_ "golang.org/x/image/webp"
)

// GuessImageFormat Guess image format from gif/jpeg/png/webp
func GuessImageFormat(data []byte) (format string) {
	reader := bytes.NewReader(data)
	_, format, err := image.DecodeConfig(reader)
	if err != nil {
		logrus.Errorf("failed to detect image %v", err)
		format = ""
	} else {
		format = "." + format
	}
	return
}
