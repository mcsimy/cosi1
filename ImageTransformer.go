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
	"math"
)

const HISTOGRAM_WIDTH = 256
const FILTER_MASK_RATIO = 0.1

type Histogram [HISTOGRAM_WIDTH]int

type ImageTransformer struct {
	sourceImage image.Image
	transformedImage image.Image
	filteredImage image.Image
	improvedImage image.Image
}

type ImageFilterInterface interface {
	 GetColor(x, y int) (result color.Color, err error)
}

type PixelTransformerInterface interface {
	 GetColor(c color.Color) (result color.Color, err error)
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

func (transformer *ImageTransformer) FilterImage(maskSize int) (err error) {
	var filter ImageFilterInterface
	filter, err = NewMedianImageFilter(transformer.sourceImage, maskSize)
	if err != nil {
		return err
	}

	resultImage := image.NewRGBA(transformer.sourceImage.Bounds())
	bounds := transformer.sourceImage.Bounds()
	var c color.Color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if c, err = filter.GetColor(x, y); err != nil {
				return err
			}
			resultImage.Set(x, y, c)
		}
	}

	transformer.filteredImage = resultImage

	return nil
}

func (transformer *ImageTransformer) TransformImage()(err error) {
	if transformer.sourceImage == nil {
		return errors.New("SOurce image not initialized")
	}

	resultImage := image.NewRGBA(transformer.sourceImage.Bounds())
	bounds := transformer.sourceImage.Bounds()
	var pixelTransformer PixelTransformerInterface = NewNegativeTransformer()
	var c color.Color

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if c, err = pixelTransformer.GetColor(transformer.sourceImage.At(x, y)); err != nil {
				return err
			}
			resultImage.Set(x, y, c)
		}
	}

	transformer.transformedImage = resultImage

	return nil
}

func (transformer *ImageTransformer) ImproveImage(brightnessPercentage, contrastPercentage int)(err error) {
	if transformer.sourceImage == nil {
		return errors.New("SOurce image not initialized")
	}

	brightnessRatio := float64(brightnessPercentage) / 100
	brightnessRatio = math.Max(-1, brightnessRatio)
	brightnessRatio = math.Min(1, brightnessRatio)

	contrastRatio := float64((100 + float64(contrastPercentage)) / 200)
	contrastRatio = math.Max(0, contrastRatio)
	contrastRatio = math.Min(1, contrastRatio)


	resultImage := image.NewRGBA(transformer.sourceImage.Bounds())
	bounds := transformer.sourceImage.Bounds()
	var pixelTransformer PixelTransformerInterface = NewBrightnessContrastTransformer(brightnessRatio, contrastRatio)
	var c color.Color

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if c, err = pixelTransformer.GetColor(transformer.sourceImage.At(x, y)); err != nil {
				return err
			}
			resultImage.Set(x, y, c)
		}
	}

	transformer.improvedImage = resultImage

	return nil
}

func (transformer *ImageTransformer) DumpSourceImage(w io.Writer) error {
	return transformer.DumpImage(transformer.sourceImage, w)
}

func (transformer *ImageTransformer) DumpTransformedImage(w io.Writer) error {
	return transformer.DumpImage(transformer.transformedImage, w)
}

func (transformer *ImageTransformer) DumpImprovedImage(w io.Writer) error {
	return transformer.DumpImage(transformer.improvedImage, w)
}

func (transformer *ImageTransformer) DumpFilteredImage(w io.Writer) error {
	return transformer.DumpImage(transformer.filteredImage, w)
}

func (transformer *ImageTransformer) DumpImage(targetImage image.Image, w io.Writer) error {
	if targetImage == nil {
		return errors.New("target image not initialized")
	}
	return png.Encode(w, targetImage);
}

func (transformer *ImageTransformer) GetSourceHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.sourceImage)
}

func (transformer *ImageTransformer) GetTransformedHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.transformedImage)
}

func (transformer *ImageTransformer) GetFilteredHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.filteredImage)
}

func (transformer *ImageTransformer) GetImprovedHistogram() (hist Histogram, err error) {
	return transformer.getHistogram(transformer.improvedImage)
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

