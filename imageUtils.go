package main

import (
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type PngFile struct {
	Image     image.Image
	ImageName string
}

func ConcatImage(picDir string) error {
	pngs, err := DecodePngDir(picDir)
	if err != nil {
		return xerrors.Errorf("faild to decode pngs: %v", err)
	}

	var pngside, pngvartical, outputvartical, outputside int
	switch {
	case len(pngs) == 0:
		return xerrors.Errorf("no image directory")
	default:
		pngside = pngs[0].Image.Bounds().Dx()
		pngvartical = pngs[0].Image.Bounds().Dy()
		outputside = pngside * 2
		if len(pngs)%2 == 1 {
			outputvartical = pngvartical * ((len(pngs) + 1) / 2)
		} else {
			outputvartical = pngvartical * (len(pngs) / 2)
		}
	}
	outputimg := image.Rectangle{image.Point{0, 0}, image.Point{outputside, outputvartical}}

	var Recpositions []image.Rectangle
	p0 := image.Point{0, 0}
	for p0.Y != outputimg.Max.Y {
		p1, p2, p3, p4, p5 := CalculateImagePoint(p0, pngside, pngvartical)
		Recpositions = append(Recpositions, image.Rectangle{p1, p2}, image.Rectangle{p3, p4})
		p0 = p5
	}

	rgba := image.NewRGBA(outputimg)
	for i := 0; i < len(pngs); i++ {
		draw.Draw(rgba, Recpositions[i], pngs[i].Image, image.Point{0, 0}, draw.Src)
	}
	outfileName := pngs[0].ImageName[:6] + "-" + pngs[len(pngs)-1].ImageName[:6]
	out, err := os.Create(outfileName + ".png")
	if err != nil {
		return xerrors.Errorf("faild to create image file: %v", err)
	}

	return png.Encode(out, rgba)
}

func CalculateImagePoint(p0 image.Point, side, vartical int) (p1, p2, p3, p4, p5 image.Point) {
	p1 = p0
	p2 = image.Point{p1.X + side, p1.Y + vartical}
	p3 = image.Point{p1.X + side, p1.Y}
	p4 = image.Point{p1.X + side*2, p1.Y + vartical}
	p5 = image.Point{0, p1.Y + vartical}
	return
}

func DecodePngDir(picDir string) (pngs []PngFile, err error) {
	files, err := ioutil.ReadDir(picDir)
	if err != nil {
		return nil, xerrors.Errorf("faild to read pictures directory: %v", err)
	}

	for _, file := range files {
		dpng, err := DecodePng(filepath.Join(picDir, file.Name()))
		if err != nil {
			return nil, xerrors.Errorf("faild to decode %v: %v", file.Name(), err)
		}
		png := new(PngFile)
		png.Image = dpng
		png.ImageName = file.Name()
		pngs = append(pngs, *png)
	}
	return
}

func DecodePng(pngname string) (img image.Image, err error) {
	f, err := os.Open(pngname)
	if err != nil {
		return nil, err
	}
	img, err = png.Decode(f)
	if err != nil {
		return nil, err
	}
	return
}
