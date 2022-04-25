package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/cenkalti/dominantcolor"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

func GetImageFromBase64(data string) (image.Image, error) {
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return nil, fmt.Errorf("invalid base64 data")
	}
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[idx+8:]))
	buff := bytes.Buffer{}
	_, err := buff.ReadFrom(reader)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(buff.Bytes()))
	return img, err
}

func GetDominantColor(image image.Image) int {
	foundColor := dominantcolor.Find(image)
	r, g, b, _ := foundColor.R, foundColor.G, foundColor.B, foundColor.A
	result := int(r)*256*256 + int(g)*256 + int(b)

	return result
}
