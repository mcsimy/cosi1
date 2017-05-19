package main

import "image/color"

type NegativeTransformer struct {

}

func NewNegativeTransformer() NegativeTransformer {
	return NegativeTransformer{}
}

func (transformer NegativeTransformer) GetColor (c color.Color) (result color.Color, err error)  {
	rgbac, _ := color.RGBAModel.Convert(c).(color.RGBA)

	rgbac.R = ^rgbac.R
	rgbac.G = ^rgbac.G
	rgbac.B = ^rgbac.B

	return rgbac, nil
}
