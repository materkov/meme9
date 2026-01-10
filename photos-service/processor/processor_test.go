package processor

import (
	"bytes"
	"context"
	"encoding/binary"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

// makeTestImage creates a simple deterministic RGBA image.
func makeTestImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Fill with a simple pattern so encoding/decoding does something non-trivial.
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{
				R: uint8((x * 3) % 255),
				G: uint8((y * 7) % 255),
				B: uint8((x + y) % 255),
				A: 255,
			})
		}
	}
	return img
}

func encodeJPEG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("jpeg.Encode: %v", err)
	}
	return buf.Bytes()
}

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func decodeJPEGSize(t *testing.T, b []byte) (int, int) {
	t.Helper()
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("jpeg.Decode: %v", err)
	}
	r := img.Bounds()
	return r.Dx(), r.Dy()
}

func TestProcess_EmptyFile(t *testing.T) {
	p := New()

	_, err := p.Process(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error for empty file, got nil")
	}
}

func TestProcess_RejectsUnsupportedType(t *testing.T) {
	p := New()

	_, err := p.Process(context.Background(), []byte("not an image"))
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
}

func TestProcess_JPEG_PreservesIfWithinMax(t *testing.T) {
	p := New()

	in := encodeJPEG(t, makeTestImage(200, 300))
	out, err := p.Process(context.Background(), in)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	w, h := decodeJPEGSize(t, out)
	if w != 200 || h != 300 {
		t.Fatalf("size: expected 200x300, got %dx%d", w, h)
	}
}

func TestProcess_JPEG_ResizesDownToMax400(t *testing.T) {
	p := New()

	// Large landscape; should downscale so width becomes 400 and height scales proportionally.
	in := encodeJPEG(t, makeTestImage(1200, 600))

	out, err := p.Process(context.Background(), in)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	w, h := decodeJPEGSize(t, out)
	if w > 400 || h > 400 {
		t.Fatalf("size: expected max 400x400, got %dx%d", w, h)
	}
	// 1200x600 => scale 400/1200=0.333.. => 400x200
	if w != 400 || h != 200 {
		t.Fatalf("size: expected 400x200, got %dx%d", w, h)
	}
}

func TestProcess_PNG_ResizesDownToMax400(t *testing.T) {
	p := New()

	// Note: current Processor.Process always encodes output as JPEG (even for PNG inputs).
	// This test validates resize behavior using a PNG input but decodes output as JPEG.
	in := encodePNG(t, makeTestImage(600, 1200))

	out, err := p.Process(context.Background(), in)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	w, h := decodeJPEGSize(t, out)
	if w > 400 || h > 400 {
		t.Fatalf("size: expected max 400x400, got %dx%d", w, h)
	}
	// 600x1200 => scale 400/1200=0.333.. => 200x400
	if w != 200 || h != 400 {
		t.Fatalf("size: expected 200x400, got %dx%d", w, h)
	}
}

// addFakeJPEGEXIF injects a minimal APP1 EXIF segment right after SOI (FFD8).
// This is not meant to be a complete EXIF implementation; it's sufficient to ensure
// the output has been re-encoded without carrying APP1 through.
func addFakeJPEGEXIF(t *testing.T, jpegBytes []byte) []byte {
	t.Helper()
	if len(jpegBytes) < 2 || jpegBytes[0] != 0xFF || jpegBytes[1] != 0xD8 {
		t.Fatalf("not a JPEG SOI")
	}

	// APP1 marker: FFE1, then big-endian length including the length field itself.
	// Payload starts with "Exif\0\0".
	payload := append([]byte("Exif\x00\x00"), []byte("FAKEEXIFDATA")...)
	seg := make([]byte, 0, 2+2+2+len(payload))
	seg = append(seg, 0xFF, 0xE1)

	// length includes these 2 length bytes + payload
	length := uint16(2 + len(payload))
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, length)
	seg = append(seg, lenBytes...)
	seg = append(seg, payload...)

	// Insert right after SOI.
	out := make([]byte, 0, len(jpegBytes)+len(seg))
	out = append(out, jpegBytes[:2]...)
	out = append(out, seg...)
	out = append(out, jpegBytes[2:]...)
	return out
}

func TestProcess_StripsJPEGEXIF_ByReencode(t *testing.T) {
	p := New()

	base := encodeJPEG(t, makeTestImage(300, 300))
	withEXIF := addFakeJPEGEXIF(t, base)

	// Sanity check: ensure the injected marker exists.
	if !bytes.Contains(withEXIF, []byte("Exif\x00\x00")) {
		t.Fatalf("expected injected EXIF marker in input")
	}

	out, err := p.Process(context.Background(), withEXIF)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}

	// Re-encoding via jpeg.Encode should not preserve APP1/Exif from input.
	if bytes.Contains(out, []byte("Exif\x00\x00")) {
		t.Fatalf("expected EXIF to be stripped from output")
	}

	// Output should still be a valid JPEG.
	if _, _, err := image.Decode(bytes.NewReader(out)); err != nil {
		t.Fatalf("output is not decodable image: %v", err)
	}
}
