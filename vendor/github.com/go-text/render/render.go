package render

import (
	"image/color"
	"image/draw"
	"math"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/opentype/api"
	"github.com/go-text/typesetting/shaping"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/math/fixed"
)

// Renderer defines a type that can render strings to a bitmap canvas.
// The size and look of output depends on the various fields in this struct.
// Developers should provide suitable output images for their draw requests.
// This type is not thread safe so instances should be used from only 1 goroutine.
type Renderer struct {
	// FontSize defines the point size of output text, commonly between 10 and 14 for regular text
	FontSize float32
	// PixScale is used to indicate the pixel density of your output target.
	// For example on a hi-DPI (or "retina") display this may be 2.0.
	// Default value is 1.0, meaning 1 pixel on the image for each render pixel.
	PixScale float32
	// Color is the pen colour for rendering
	Color color.Color

	shaper      shaping.Shaper
	filler      *rasterx.Filler
	fillerScale float32
}

// DrawString will rasterise the given string into the output image using the specified font face.
// The text will be drawn starting at the left edge, down from the image top by the
// font ascent value, so that the text is all visible.
// The return value is the X pixel position of the end of the drawn string.
func (r *Renderer) DrawString(str string, img draw.Image, face font.Face) int {
	if r.PixScale == 0 {
		r.PixScale = 1
	}

	in := shaping.Input{
		Text:     []rune(str),
		RunStart: 0,
		RunEnd:   len(str),
		Face:     face,
		Size:     fixed.I(int(r.FontSize)),
	}
	out := r.cachedShaper().Shape(in)
	return r.DrawShapedRunAt(out, img, 0, out.LineBounds.Ascent.Ceil())
}

// DrawStringAt will rasterise the given string into the output image using the specified font face.
// The text will be drawn starting at the x, y pixel position.
// Note that x and y are not multiplied by the `PixScale` value as they refer to output coordinates.
// The return value is the X pixel position of the end of the drawn string.
func (r *Renderer) DrawStringAt(str string, img draw.Image, x, y int, face font.Face) int {
	if r.PixScale == 0 {
		r.PixScale = 1
	}

	in := shaping.Input{
		Text:     []rune(str),
		RunStart: 0,
		RunEnd:   len(str),
		Face:     face,
		Size:     fixed.I(int(r.FontSize)),
	}
	return r.DrawShapedRunAt(r.cachedShaper().Shape(in), img, x, y)
}

// DrawShapedRunAt will rasterise the given shaper run into the output image using font face referenced in the shaping.
// The text will be drawn starting at the startX, startY pixel position.
// Note that startX and startY are not multiplied by the `PixScale` value as they refer to output coordinates.
// The return value is the X pixel position of the end of the drawn string.
func (r *Renderer) DrawShapedRunAt(run shaping.Output, img draw.Image, startX, startY int) int {
	if r.PixScale == 0 {
		r.PixScale = 1
	}
	scale := r.FontSize * r.PixScale / float32(run.Face.Upem())
	r.fillerScale = scale

	b := img.Bounds()
	scanner := rasterx.NewScannerGV(b.Dx(), b.Dy(), img, b)
	f := rasterx.NewFiller(b.Dx(), b.Dy(), scanner)
	r.filler = f
	f.SetColor(r.Color)
	x := float32(startX)
	y := float32(startY)
	for _, g := range run.Glyphs {
		xPos := x + fixed266ToFloat(g.XOffset)*r.PixScale
		yPos := y - fixed266ToFloat(g.YOffset)*r.PixScale
		data := run.Face.GlyphData(g.GlyphID)
		switch format := data.(type) {
		case api.GlyphOutline:
			r.drawOutline(g, format, f, scale, xPos, yPos)
		case api.GlyphBitmap:
			_ = r.drawBitmap(g, format, img, xPos, yPos)
		case api.GlyphSVG:
			_ = r.drawSVG(g, format, img, xPos, yPos)
		}

		x += fixed266ToFloat(g.XAdvance) * r.PixScale
	}
	f.Draw()
	r.filler = nil
	return int(math.Ceil(float64(x)))
}

func (r *Renderer) cachedShaper() shaping.Shaper {
	if r.shaper == nil {
		r.shaper = &shaping.HarfbuzzShaper{}
	}

	return r.shaper
}

func (r *Renderer) drawOutline(g shaping.Glyph, bitmap api.GlyphOutline, f *rasterx.Filler, scale float32, x, y float32) {
	for _, s := range bitmap.Segments {
		switch s.Op {
		case api.SegmentOpMoveTo:
			f.Start(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)})
		case api.SegmentOpLineTo:
			f.Line(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)})
		case api.SegmentOpQuadTo:
			f.QuadBezier(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)},
				fixed.Point26_6{X: floatToFixed266(s.Args[1].X*scale + x), Y: floatToFixed266(-s.Args[1].Y*scale + y)})
		case api.SegmentOpCubeTo:
			f.CubeBezier(fixed.Point26_6{X: floatToFixed266(s.Args[0].X*scale + x), Y: floatToFixed266(-s.Args[0].Y*scale + y)},
				fixed.Point26_6{X: floatToFixed266(s.Args[1].X*scale + x), Y: floatToFixed266(-s.Args[1].Y*scale + y)},
				fixed.Point26_6{X: floatToFixed266(s.Args[2].X*scale + x), Y: floatToFixed266(-s.Args[2].Y*scale + y)})
		}
	}
	f.Stop(true)
}

func fixed266ToFloat(i fixed.Int26_6) float32 {
	return float32(float64(i) / 64)
}

func floatToFixed266(f float32) fixed.Int26_6 {
	return fixed.Int26_6(int(float64(f) * 64))
}
