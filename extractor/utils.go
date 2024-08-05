/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package extractor

import (
	"fmt"
	"image/color"

	"github.com/showntop/unipdf/common"
	"github.com/showntop/unipdf/core"
	"github.com/showntop/unipdf/model"
)

// RenderMode specifies the text rendering mode (Tmode), which determines whether showing text shall cause
// glyph outlines to be  stroked, filled, used as a clipping boundary, or some combination of the three.
// Stroking, filling, and clipping shall have the same effects for a text object as they do for a path object
// (see 8.5.3, "Path-Painting Operators" and 8.5.4, "Clipping Path Operators").
type RenderMode int

// Render mode type.
const (
	RenderModeStroke RenderMode = 1 << iota // Stroke
	RenderModeFill                          // Fill
	RenderModeClip                          // Clip
)

// toFloatXY returns `objs` as 2 floats, if that's what `objs` is, or an error if it isn't.
func toFloatXY(objs []core.PdfObject) (x, y float64, err error) {
	if len(objs) != 2 {
		return 0, 0, fmt.Errorf("invalid number of params: %d", len(objs))
	}
	floats, err := core.GetNumbersAsFloat(objs)
	if err != nil {
		return 0, 0, err
	}
	return floats[0], floats[1], nil
}

// truncate returns the first `n` characters in string `s`.
func truncate(s string, n int) string {
	if len(s) < n {
		return s
	}
	return s[:n]
}

// pdfColorToGoColor converts the specified color to a Go color, using the
// provided colorspace. If unsuccessful, color.Black is returned.
func pdfColorToGoColor(space model.PdfColorspace, c model.PdfColor) color.Color {
	if space == nil || c == nil {
		return color.Black
	}

	conv, err := space.ColorToRGB(c)
	if err != nil {
		common.Log.Debug("WARN: could not convert color %v (%v) to RGB: %s", c, space, err)
		return color.Black
	}
	rgb, ok := conv.(*model.PdfColorDeviceRGB)
	if !ok {
		common.Log.Debug("WARN: converted color is not in the RGB colorspace: %v", conv)
		return color.Black
	}

	return color.NRGBA{
		R: uint8(rgb.R() * 255),
		G: uint8(rgb.G() * 255),
		B: uint8(rgb.B() * 255),
		A: uint8(255),
	}
}
