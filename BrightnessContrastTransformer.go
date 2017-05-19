package main

import (
	"image/color"
	"math"
	"fmt"
)

const UINT8_MEDIAN = math.MaxUint8 / 2

type BrightnessContrastTransformer struct {
	brightnessNatural float64
	contrastNatural float64
}

func NewBrightnessContrastTransformer(brightnessRatio, contrastRatio float64) BrightnessContrastTransformer {
	brightnessNatural := brightnessRatio * math.MaxUint8
	contrastNatural := math.Tan(contrastRatio * math.Pi / 2)
	fmt.Printf("B: %v, C: %v\n", brightnessNatural, contrastNatural)

	return BrightnessContrastTransformer{brightnessNatural: brightnessNatural, contrastNatural: contrastNatural}
}

func (transformer BrightnessContrastTransformer) GetColor (c color.Color) (result color.Color, err error)  {
	rgbac, _ := color.RGBAModel.Convert(c).(color.RGBA)

	rgbac.R = transformer.applyContrast(rgbac.R)
	rgbac.R = transformer.applyBrightness(rgbac.R)
	rgbac.G = transformer.applyContrast(rgbac.G)
	rgbac.G = transformer.applyBrightness(rgbac.G)
	rgbac.B = transformer.applyContrast(rgbac.B)
	rgbac.B = transformer.applyBrightness(rgbac.B)

	return rgbac, nil
}

func (transformer BrightnessContrastTransformer)applyContrast(channelValue uint8) uint8 {
	result := transformer.contrastNatural * (float64(channelValue) - UINT8_MEDIAN) + UINT8_MEDIAN
	result = math.Max(result, 0)
	result = math.Min(result, math.MaxUint8)
	return uint8(result)
}

func (transformer BrightnessContrastTransformer)applyBrightness(channelValue uint8) uint8 {
	result := transformer.brightnessNatural + float64(channelValue)
	result = math.Max(result, 0)
	result = math.Min(result, math.MaxUint8)
	return uint8(result)
}
