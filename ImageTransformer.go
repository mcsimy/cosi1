package main

import (
	"image"
	"io"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"errors"
	"image/png"
	"image/color"
)

const HISTOGRAM_WIDTH = 256

type Histogram [HISTOGRAM_WIDTH]int

type ImageTransformer struct {
	sourceImage image.Image
	transformedImage image.Image
	filteredImage image.Image
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

func (transformer *ImageTransformer) TransformImage() error {
	if transformer.sourceImage == nil {
		return errors.New("SOurce image not initialized")
	}

	var err error
	transformer.transformedImage, err = transformer.getNegativeImage(transformer.sourceImage)

	return err
}

func (transformer *ImageTransformer) getNegativeImage(sourceImage image.Image) (image.Image, error) {
	if sourceImage == nil {
		return nil, errors.New("SourceImage not initialized")
	}

	resultImage := image.NewRGBA(transformer.sourceImage.Bounds())

	bounds := sourceImage.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c, _ := color.RGBAModel.Convert(sourceImage.At(x, y)).(color.RGBA)

			c.R = ^c.R
			c.G = ^c.G
			c.B = ^c.B
			resultImage.Set(x, y, c)
		}
	}

	return resultImage, nil
}


func (transformer *ImageTransformer) DumpSourceImage(w io.Writer) error {
	if transformer.sourceImage == nil {
		return errors.New("Source image not initialized")
	}
	return png.Encode(w, transformer.sourceImage);
}

func (transformer *ImageTransformer) DumpTransformedImage(w io.Writer) error {
	if transformer.transformedImage == nil {
		return errors.New("Transformed image not initialized")
	}
	return png.Encode(w, transformer.transformedImage);
}

func (transformer *ImageTransformer) GetSourceHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.sourceImage)
}

func (transformer *ImageTransformer) GetTransformedHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.transformedImage)
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

