package steganography

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

type WorkParams struct {
	InputPath      string
	OutputPath     string
	DataString     string
	Channels       string
	BitsPerChannel int
}

func loadImage(filepath string) image.Image {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured whilst opening file '%s':\n%v\n", filepath, err)
		os.Exit(1)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "an error occured whilst decoding image '%s':\n%v\n", filepath, err)
		os.Exit(1)
	}

	return img
}

func Size(data string) int {
	return len([]byte(data))*8 + 32
}

func Capacity(filepath string, bpc int, channels string) int {
	img := loadImage(filepath)
	return capacity(img.Bounds(), bpc, len(channels))
}

func capacity(rect image.Rectangle, bpc int, channels int) int {
	return rect.Dx() * rect.Dy() * channels * bpc
}

func Extract(params WorkParams) string {
	img := loadImage(params.InputPath)
	extracted := decode(params, img)
	return string(extracted)
}

func Embed(params WorkParams) {
	img := loadImage(params.InputPath)

	newImg := encode(params, img)

	fi, err := os.Create(fmt.Sprintf("./%s", params.OutputPath))
	if err != nil {
		fmt.Printf("Could not create file '%s.png'.\n%v", params.OutputPath, err)
		fmt.Scanf("h")
		os.Exit(1)
	}
	defer fi.Close()

	err = png.Encode(fi, newImg)
	if err != nil {
		fmt.Printf("Failed saving image to file.\n%v", err)
		fmt.Scanf("h")
		os.Exit(1)
	}

}

func encode(params WorkParams, src image.Image) image.Image {
	data := []byte(params.DataString)
	rect := src.Bounds()
	cap := capacity(rect, params.BitsPerChannel, len(params.Channels))
	payloadSize := Size(params.DataString)

	if payloadSize > cap {
		fmt.Fprintf(os.Stderr, "insufficient capacity(%d) for payload(%d)\n", cap, payloadSize)
		os.Exit(1)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	payload := append(header, data...)

	bits := make([]byte, payloadSize)
	for i := rect.Min.X; i < len(payload); i++ {
		for j := 0; j < 8; j++ {
			index := (i*8 + j)
			bits[index] = (payload[i] >> j) & 1
		}
	}

	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), src, rect.Min, draw.Src)

	bpc := params.BitsPerChannel

	chnlCount := len(params.Channels)
	chnlIndices := make([]int, chnlCount)
	for i := 0; i < chnlCount; i++ {
		switch string(params.Channels[i]) {
		case "r":
			chnlIndices[i] = 0
		case "g":
			chnlIndices[i] = 1
		case "b":
			chnlIndices[i] = 2
		case "a":
			chnlIndices[i] = 3
		}
	}

	dataChnlCount := payloadSize / bpc
	if payloadSize%bpc != 0 {
		dataChnlCount++
	}

	var mask byte = 255 & (255 << bpc)
	chnlData := make([]*byte, dataChnlCount)
	for i := range chnlData {
		index := (i/chnlCount)*4 + chnlIndices[i%chnlCount]
		chnlData[i] = &img.Pix[index]
		*chnlData[i] = (*chnlData[i]) & mask
	}

	for i := range bits {
		bit := *chnlData[i/bpc] | (bits[i] << (i % bpc))
		*chnlData[i/bpc] = bit
	}

	return img
}

func decodeBytes(bits []byte) []byte {

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

func extractBitsFromImageChannels(img *image.RGBA, amount int, chnlIndices []int, params WorkParams) []byte {
	bpc := params.BitsPerChannel
	dataChnlCount := amount/bpc + amount%bpc
	chnlData := make([]*byte, dataChnlCount)
	chnlCount := len(params.Channels)

	for i := range chnlData {
		index := (i/chnlCount)*4 + chnlIndices[i%chnlCount]
		chnlData[i] = &img.Pix[index]
	}

	buffer := make([]byte, amount)
	for i := range buffer {
		buffer[i] = buffer[i] | ((*chnlData[i/bpc] >> (i % bpc)) & 1)
	}

	return buffer
}

func decode(params WorkParams, src image.Image) []byte {
	rect := src.Bounds()
	chnlCount := len(params.Channels)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), src, rect.Min, draw.Src)

	chnlIndices := make([]int, chnlCount)
	for i := 0; i < chnlCount; i++ {
		switch string(params.Channels[i]) {
		case "r":
			chnlIndices[i] = 0
		case "g":
			chnlIndices[i] = 1
		case "b":
			chnlIndices[i] = 2
		case "a":
			chnlIndices[i] = 3
		}
	}

	headerBits := extractBitsFromImageChannels(img, 32, chnlIndices, params)
	headerBytes := decodeBytes(headerBits)
	header := binary.BigEndian.Uint32(headerBytes)

	dataBits := extractBitsFromImageChannels(img, int(header)*8+32, chnlIndices, params)
	dataBytes := decodeBytes(dataBits[32:])
	return dataBytes
}
