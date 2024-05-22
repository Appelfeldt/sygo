package steganography

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"
)

func LoadImage(filepath string) (image.Image, error) {
	f, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func Size(filepath string) int {
	img, err := LoadImage(filepath)
	if err != nil {
		fmt.Printf("Error loading image:\n %s", err)
		os.Exit(1)
	}

	return size(img.Bounds())
}

func Extract(filepath string) string {
	img, err := LoadImage(filepath)
	if err != nil {
		fmt.Printf("Error loading image:\n %s", err)
		os.Exit(1)
	}

	extracted := extract(img)
	return string(extracted)
}

func Embed(args ...string) {
	img, err := LoadImage(args[0])
	if err != nil {
		fmt.Printf("Error loading image:\n %s", err)
		os.Exit(1)
	}

	data := []byte(args[1])

	newImg, err := insert(data, img)
	if err != nil {
		fmt.Printf("Error inserting payload:\n%s", err)
		fmt.Scanf("h")
		os.Exit(1)
	}

	outpath := args[2]
	if ext := path.Ext(outpath); ext == "" {
		outpath += ".png"
	}

	fi, err := os.Create(fmt.Sprintf("./%s", outpath))
	if err != nil {
		fmt.Printf("Could not create file '%s.png'.\n%s", outpath, err)
		fmt.Scanf("h")
		os.Exit(1)
	}
	defer fi.Close()

	err = png.Encode(fi, newImg)
	if err != nil {
		fmt.Printf("Failed saving image to file.\n%s", err)
		fmt.Scanf("h")
		os.Exit(1)
	}

}

func size(rect image.Rectangle) int {
	return rect.Dx() * rect.Dy() * 3
}

func insert(data []byte, src image.Image) (image.Image, error) {
	rect := src.Bounds()
	bitspace := rect.Dx() * rect.Dy() * 3
	bitcount := len(data)*8 + 64*8 //64*8 extra needed for uint64 header

	if bitcount > bitspace {
		return nil, errors.New("image too small for data payload")
	}

	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, uint64(len(data)))
	payload := append(header, data...)

	bits := make([]byte, bitcount)
	for i := rect.Min.X; i < len(payload); i++ {
		for j := rect.Min.Y; j < 8; j++ {
			index := (i*8 + j)
			bits[index] = (payload[i] >> j) & 1
		}
	}

	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), src, rect.Min, draw.Src)

	p := 0
	for i := 0; i < bitcount; i++ {
		if (p+1)%4 == 0 {
			p++
		}
		img.Pix[p] = img.Pix[p] | bits[i]
		p++
	}

	return img, nil
}

func extractBytes(bits []byte) []byte {

	bytecount := len(bits) / 8
	bytebuffer := make([]byte, bytecount)
	for i := 0; i < bytecount; i++ {
		b := bits[i*8 : i*8+8]
		var buffer byte = 0
		for j := 0; j < len(b); j++ {
			buffer += (byte)((b[j] & 1) << j)
		}
		bytebuffer[i] = buffer
	}
	return bytebuffer
}

func extract(src image.Image) []byte {
	rect := src.Bounds()

	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), src, rect.Min, draw.Src)
	pixels := len(img.Pix) / 4
	bitcount := 3 * pixels
	payload := make([]byte, bitcount)

	for i := 0; i < pixels; i++ {
		payload[i*3] = img.Pix[i*4]
		payload[i*3+1] = img.Pix[i*4+1]
		payload[i*3+2] = img.Pix[i*4+2]
	}

	headerbits := payload[:64]
	header := extractBytes(headerbits)
	bytecount := binary.BigEndian.Uint64(header)

	databits := payload[64 : 64+bytecount*8]
	data := extractBytes(databits)

	return data
}
