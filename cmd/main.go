package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func main() {
	text := "Ashleigh & Eamon"
	// text := "Albert"

	// file, err := os.OpenFile("fonts/MysteryQuest-Regular.ttf", os.O_RDONLY, 0644)
	file, err := os.OpenFile("fonts/FleurDeLeah-Regular.ttf", os.O_RDONLY, 0644)

	// file, err := os.OpenFile("fonts/MysteryQuest-Regular.ttf", os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}

	loader, err := opentype.NewLoader(file)
	if err != nil {
		fmt.Println(err)
	}
	f, err := font.NewFont(loader)
	if err != nil {
		fmt.Println(err)
	}
	face := font.NewFace(f)

	var advance float32 = 0

	c := canvas.New(10000, 1000)

	ctx := canvas.NewContext(c)
	// ctx.Scale(0.5, 0.5)
	ctx.SetFillColor(color.RGBA{R: 255, G: 255, B: 255, A: 0})
	// ctx.SetFillColor(color.RGBA{R: 0, G: 255, B: 0, A: 0})
	ctx.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	ctx.SetStrokeWidth(1)
	x := 0.0

	union := &canvas.Path{}
	for _, rune := range text {
		gid, ok := face.NominalGlyph(rune)

		if !ok {
			fmt.Println("no glyph")
			return
		}
		gd := face.GlyphData(gid)
		advance = face.HorizontalAdvance(gid)
		// fmt.Println(advance)
		outline, ok := gd.(font.GlyphOutline)
		if !ok {
			fmt.Println("not outline")
			return
		}

		path := &canvas.Path{}

		isOpen := false
		for _, contour := range outline.Segments {
			// fmt.Println(contour)
			switch contour.Op {
			case opentype.SegmentOpMoveTo:
				if isOpen {
					path.Close()
					isOpen = false

				}
				path.MoveTo(float64(contour.Args[0].X), float64(contour.Args[0].Y))
				isOpen = true

			case opentype.SegmentOpLineTo:
				path.LineTo(float64(contour.Args[0].X), float64(contour.Args[0].Y))
			case opentype.SegmentOpQuadTo:
				path.QuadTo(float64(contour.Args[0].X), float64(contour.Args[0].Y), float64(contour.Args[1].X), float64(contour.Args[1].Y))
			case opentype.SegmentOpCubeTo:
				fmt.Println("have a cubic")
				path.CubeTo(float64(contour.Args[0].X), float64(contour.Args[0].Y), float64(contour.Args[1].X), float64(contour.Args[1].Y), float64(contour.Args[2].X), float64(contour.Args[2].Y))
			}
		}
		path.Close()
		path = path.Settle(canvas.EvenOdd)
		path.Translate(x, 0)
		ctx.DrawPath(10, 10, path)
		// ctx.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		// ctx.DrawPath(10, 10, path)
		// ctx.SetStrokeColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
		x += float64(advance)
		for range 20 {
			path = path.Offset(-5, 0.1)
			path = path.Settle(canvas.EvenOdd)
			path = path.SimplifyVisvalingamWhyatt(1)
			ctx.DrawPath(10, 10, path)
		}

		union = union.Or(path)

		// ctx.Translate(float64(advance), 0)
	}
	// fmt.Println(union)
	// union = union.SimplifyVisvalingamWhyatt(1)
	// fmt.Println(union.Len())
	// union = union.Settle(canvas.EvenOdd)
	// ctx.DrawPath(10, 10, union)
	// ctx.SetFillColor(color.RGBA{R: 0, G: 255, B: 0, A: 0})
	// ctx.SetStrokeColor(color.RGBA{R: 0, G: 255, B: 0, A: 255})
	// ps := []*canvas.Path{union}
	//
	// for i := 1; i < 20; i++ {
	// 	next := ps[i-1].Offset(-3, 1)
	// 	next = next.SimplifyVisvalingamWhyatt(1)
	// 	next = next.Settle(canvas.EvenOdd)
	// 	fmt.Println(i, next.Len())
	// 	ctx.DrawPath(10, 10, next)
	// 	ps = append(ps, next)
	// }
	// ctx.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 255, A: 255})
	// ctx.SetStrokeWidth(2)
	// ctx.DrawPath(10, 10, union)

	c.Fit(20)
	err = renderers.Write("test.png", c, canvas.DPMM(3.2))
	if err != nil {
		fmt.Println(err)
	}

}
