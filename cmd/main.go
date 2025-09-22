package main

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

const (
	FEED      = 2.0
	THICKNESS = 10.0
	STEPOVER  = 0.5
	SAFE_Z    = THICKNESS + 2
	LENGTH    = 485.0
	WIDTH     = 890.0
)

type Datum struct {
	Text  string
	Font  string
	Ox    float64
	Oy    float64
	Scale float64
}

func main() {
	xMargin := 40.0
	// nameFont := "Zeyada-Regular.ttf"
	// dateFont := "BaiJamjuree-Regular.ttf"
	// data := []Datum{
	// 	{
	// 		Text:  "Ashleigh",
	// 		Font:  nameFont,
	// 		Ox:    xMargin,
	// 		Oy:    700.0,
	// 		Scale: 0.14,
	// 	},
	// 	{
	// 		Text:  "& Eamon",
	// 		Font:  nameFont,
	// 		Ox:    25.0,
	// 		Oy:    550.0,
	// 		Scale: 0.14,
	// 	},
	// 	{
	// 		Text: "26·10·2025",
	//
	// 		Font:  dateFont,
	// 		Ox:    126.0,
	// 		Oy:    200.0,
	// 		Scale: 0.05,
	// 	},
	// }
	// nameFont := "LiuJianMaoCao-Regular.ttf"
	// dateFont := "Raleway-Regular.ttf"
	// data := []Datum{
	// 	{
	// 		Text:  "Ashleigh",
	// 		Font:  nameFont,
	// 		Ox:    xMargin,
	// 		Oy:    700.0,
	// 		Scale: 0.14,
	// 	},
	// 	{
	// 		Text:  "& Eamon",
	// 		Font:  nameFont,
	// 		Ox:    25.0,
	// 		Oy:    550.0,
	// 		Scale: 0.14,
	// 	},
	// 	{
	// 		Text: "26·10·2025",
	//
	// 		Font:  dateFont,
	// 		Ox:    126.0,
	// 		Oy:    200.0,
	// 		Scale: 0.05,
	// 	},
	// }
	nameFont := "MysteryQuest-Regular.ttf"
	dateFont := "Creepster-Regular.ttf"
	data := []Datum{
		{
			Text:  "Ashleigh",
			Font:  nameFont,
			Ox:    xMargin,
			Oy:    700.0,
			Scale: 0.12,
		},
		{
			Text:  "& Eamon",
			Font:  "Kablammo-Regular-VariableFont_MORF.ttf",
			Ox:    25.0,
			Oy:    550.0,
			Scale: 0.054,
		},
		{
			Text: "26·10·2025",

			Font:  dateFont,
			Ox:    126.0,
			Oy:    200.0,
			Scale: 0.05,
		},
	}
	c := canvas.New(1, 1)
	ctx := canvas.NewContext(c)
	ctx.SetStrokeWidth(1)
	ctx.SetFillColor(color.RGBA{R: 255, G: 255, B: 255, A: 0})
	ctx.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	union := &canvas.Path{}
	file, err := os.Create("test.nc")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for _, datum := range data {
		x := 0.0
		fontFile, err := os.OpenFile(filepath.Join("fonts", datum.Font), os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println(err)
		}

		loader, err := opentype.NewLoader(fontFile)
		if err != nil {
			fmt.Println(err)
		}
		f, err := font.NewFont(loader)
		if err != nil {
			fmt.Println(err)
		}
		face := font.NewFace(f)

		var advance float32 = 0

		for _, rune := range datum.Text {
			gid, ok := face.NominalGlyph(rune)

			if !ok {
				fmt.Println("no glyph")
				return
			}
			gd := face.GlyphData(gid)
			advance = face.HorizontalAdvance(gid)
			outline, ok := gd.(font.GlyphOutline)
			if !ok {
				fmt.Println("not outline")
				return
			}

			path := &canvas.Path{}

			isOpen := false
			for _, contour := range outline.Segments {
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
					path.CubeTo(float64(contour.Args[0].X), float64(contour.Args[0].Y), float64(contour.Args[1].X), float64(contour.Args[1].Y), float64(contour.Args[2].X), float64(contour.Args[2].Y))
				}
			}
			ctx.DrawPath(0, 0, canvas.Rectangle(LENGTH, WIDTH))
			path.Close()
			path = path.Settle(canvas.EvenOdd)
			path = path.SimplifyVisvalingamWhyatt(1)
			path.Translate(x, 0)
			path.Scale(datum.Scale, datum.Scale)
			ctx.SetStrokeColor(color.RGBA{R: 255, G: 0, B: 255, A: 255})

			z := THICKNESS - 0.1
			file.WriteString(pathToGcode(path, 0, 0, 0.1, z))

			ctx.DrawPath(datum.Ox, datum.Oy, path)
			ctx.SetStrokeColor(color.RGBA{R: 0, G: 255, B: 255, A: 255})

			x += float64(advance)
			for range 40 {
				z -= STEPOVER
				path = path.Offset(-STEPOVER, 0.2)
				path = path.Settle(canvas.EvenOdd)
				path = path.SimplifyVisvalingamWhyatt(1)
				ctx.DrawPath(datum.Ox, datum.Oy, path)
				file.WriteString(pathToGcode(path, 0, 0, 0.1, z))
			}
			union = union.Or(path)
		}
	}

	c.Fit(20)
	err = renderers.Write("image.png", c, canvas.DPMM(3.2))
	if err != nil {
		fmt.Println(err)
	}

}

func pathToGcode(path *canvas.Path, offestX, offsetY, scale, z float64) string {
	gcode := strings.Builder{}
	subPaths := path.Split()
	for _, subPath := range subPaths {
		scanner := subPath.Scanner()
		for scanner.Scan() {
			vs := scanner.Values()
			x := vs[0]
			y := vs[1]
			x = x*scale + offestX
			y = y*scale + offsetY
			g := 0.0
			zz := SAFE_Z
			if scanner.Cmd() == 2 {
				g = 1.0
				zz = z
			}
			gcode.WriteString(gcodeHelper(g, FEED, x, y, zz))
		}
	}
	return gcode.String()
}

func gcodeHelper(g, f, x, y, z float64) string {
	return fmt.Sprintf("G%.0f F%.0f X%.3f Y%.3f Z%.3f\n", g, f, x, y, z)
}
