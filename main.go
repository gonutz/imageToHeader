package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: 1. image file path 2. header file path")
		return
	}

	imageFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer imageFile.Close()

	img, err := png.Decode(imageFile)
	if err != nil {
		panic(err)
	}

	_, imageFileName := filepath.Split(os.Args[1])
	imageName := strings.TrimSuffix(imageFileName, filepath.Ext(imageFileName))

	buffer := bytes.NewBuffer(nil)
	output := panickyWriter{buffer}

	output.WriteLine("// generated by imageToHeader.go")
	output.WriteLine("// data format is ARGB in range [0..255]")
	output.WriteLine("// the ordering is top to bottom")
	output.WriteLine(fmt.Sprintf("const unsigned int %vImageWidth = %v;",
		imageName, img.Bounds().Dx()))
	output.WriteLine(fmt.Sprintf("const unsigned int %vImageHeight = %v;",
		imageName, img.Bounds().Dy()))
	output.WriteLine(fmt.Sprintf("unsigned char %vImageData[] = {", imageName))
	colorModel := color.NRGBAModel
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		output.WriteString("   ")
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			nrgba := colorModel.Convert(img.At(x, y)).(color.NRGBA)
			output.WriteString(fmt.Sprintf(" %v, %v, %v, %v,",
				nrgba.A, nrgba.R, nrgba.G, nrgba.B))
		}
		output.WriteLine("")
	}
	buffer.Truncate(len(buffer.Bytes()) - 2)
	output.WriteLine("")
	output.WriteLine("};")

	if err := ioutil.WriteFile(os.Args[2], buffer.Bytes(), 0666); err != nil {
		panic(err)
	}
}

type stringWriter interface {
	WriteString(string) (int, error)
}

type panickyWriter struct {
	stringWriter
}

func (w panickyWriter) WriteString(s string) {
	_, err := w.stringWriter.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func (w panickyWriter) WriteLine(line string) {
	w.WriteString(line + "\n")
}
