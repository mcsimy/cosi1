package main

import (
	"image"
	"io"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"errors"
)

const HISTOGRAM_WIDTH = 256

type Histogram [HISTOGRAM_WIDTH]int

type ImageTransformer struct {
	sourceImage image.Image
}

func (transformer *ImageTransformer) LoadSourceImage(r io.Reader) error {
	var err error
	var imageFormat string

	transformer.sourceImage, imageFormat, err = image.Decode(r)

	if err != nil {
		return err
	}

	log.Printf("image format: %s", imageFormat)
	return nil
}

func (transformer *ImageTransformer) GetSourceHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.sourceImage)
}

func (transformer *ImageTransformer) getHistogram(targetImage image.Image) (hist Histogram, err error) {
	if targetImage == nil {
		return hist, errors.New("Image is not initialized")
	}
	//multipy by 3 due to max of 3 channels is 2^16 * 3
	histogramRatio := uint32(1 << 16 / HISTOGRAM_WIDTH * 3)

	bounds := targetImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := targetImage.At(x, y).RGBA()
			brightness := int((r + g + b) / histogramRatio)
			hist[brightness]++
		}
	}

	return hist, nil
}

