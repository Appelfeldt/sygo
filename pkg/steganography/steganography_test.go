package steganography

import (
	"image"
	"testing"
)

func TestEncodingAndDecoding(t *testing.T) {
	params := WorkParams{
		DataString:     "τhe quicк βrowи fθx jumped over tнe lazy dog",
		Channels:       "rga",
		BitsPerChannel: 3,
	}

	rect := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: 100,
			Y: 100,
		},
	}
	img := image.NewRGBA(rect)
	pixels := rect.Dx() * rect.Dy()
	for i := 0; i < pixels; i++ {
		img.Pix[i*4] = 0
		img.Pix[i*4+1] = 96
		img.Pix[i*4+2] = 128
		img.Pix[i*4+3] = 255
	}

	encodedImg := encode(params, img)
	decodeResult := decode(params, encodedImg)
	decodedString := string(decodeResult)
	if decodedString != params.DataString {
		t.Errorf("Expected value '%s', got '%s'\n", params.DataString, decodedString)
	}
}
