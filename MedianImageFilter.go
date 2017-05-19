package main

import (
	"math"
	"image/color"
	"sort"
	"image"
	"errors"
)

type MedianBuffer struct {
	R []int
	G []int
	B []int
}

type MedianImageFilter struct {
	sourceImage image.Image
	maskEdge int
	bounds image.Rectangle
	buffer MedianBuffer
}

func NewMedianImageFilter(sourceImage image.Image, maskSize int) (MedianImageFilter, error) {
	if sourceImage == nil {
		return MedianImageFilter{}, errors.New("SOurce image not initialized")
	}
	bounds := sourceImage.Bounds()
	sourceImageSize := math.Min(float64(bounds.Max.X - bounds.Min.X), float64(bounds.Max.Y - bounds.Min.Y))

	if maskSize % 2 == 0 || maskSize < 3 || float64(maskSize) > sourceImageSize * FILTER_MASK_RATIO {
		return MedianImageFilter{}, errors.New("Mask size should be odd number between 3 and 0.01 * MIN_IMG_SIZE")
	}

	buffer := MedianBuffer{
		make([]int, maskSize * maskSize),
		make([]int, maskSize * maskSize),
		make([]int, maskSize * maskSize),
	}

	return MedianImageFilter{sourceImage, int(maskSize / 2), bounds, buffer}, nil
}

func (filter MedianImageFilter) GetColor(pointX, pointY int) (result color.Color, err error) {
	//TODO
	i := 0

	for y := pointY - filter.maskEdge; y <= pointY + filter.maskEdge; y++ {
		for x := pointX - filter.maskEdge; x <= pointX + filter.maskEdge; x++ {
			//fix mask near image boundaries
			realX := int(math.Max(float64(x), float64(filter.bounds.Min.X)))
			realX = int(math.Min(float64(realX), float64(filter.bounds.Max.X - 1)))
			realY := int(math.Max(float64(y), float64(filter.bounds.Min.Y)))
			realY = int(math.Min(float64(realY), float64(filter.bounds.Max.Y - 1)))

			c, _ := color.RGBAModel.Convert(filter.sourceImage.At(realX, realY)).(color.RGBA)
			filter.buffer.R[i] = int(c.R)
			filter.buffer.G[i] = int(c.G)
			filter.buffer.B[i] = int(c.B)
			i++

			//fmt.Printf("%d*%d (%v)    ", realX, realY, c)
		}
		//fmt.Println("")
	}

	c, _ := color.RGBAModel.Convert(filter.sourceImage.At(pointX, pointY)).(color.RGBA)

	sort.Ints(filter.buffer.R)
	sort.Ints(filter.buffer.G)
	sort.Ints(filter.buffer.B)

	//fmt.Printf("\n%v\n", filter.buffer)

	medianIndex := int(len(filter.buffer.R) / 2) + 1

	result = color.RGBA{
		R: uint8(filter.buffer.R[medianIndex]),
		G: uint8(filter.buffer.G[medianIndex]),
		B: uint8(filter.buffer.B[medianIndex]),
		A: c.A,
	}

	return result, nil
}

