// Package mp3 provides MP3 audio decoding for Future Render.
//
// The decoder wraps github.com/hajimehoshi/go-mp3 and produces signed
// 16-bit little-endian stereo PCM data suitable for audio.Context.NewPlayer.
package mp3

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	gomp3 "github.com/hajimehoshi/go-mp3"
)

// Stream is a decoded MP3 audio stream. It provides signed 16-bit
// LE stereo PCM data.
type Stream struct {
	data       *bytes.Reader
	raw        []byte
	sampleRate int
	length     int64
}

// Read implements io.Reader.
func (s *Stream) Read(p []byte) (int, error) {
	return s.data.Read(p)
}

// Seek implements io.Seeker.
func (s *Stream) Seek(offset int64, whence int) (int64, error) {
	return s.data.Seek(offset, whence)
}

// Length returns the total length of decoded audio data in bytes.
func (s *Stream) Length() int64 {
	return s.length
}

// SampleRate returns the sample rate of the MP3 file.
func (s *Stream) SampleRate() int {
	return s.sampleRate
}

// Decode reads an MP3 file from src and returns a Stream of stereo
// 16-bit signed LE PCM data at the file's native sample rate.
func Decode(src io.Reader) (*Stream, error) {
	decoder, err := gomp3.NewDecoder(src)
	if err != nil {
		return nil, fmt.Errorf("mp3: open decoder: %w", err)
	}

	sampleRate := decoder.SampleRate()

	pcm, err := readAll(decoder)
	if err != nil {
		return nil, err
	}

	return &Stream{
		data:       bytes.NewReader(pcm),
		raw:        pcm,
		sampleRate: sampleRate,
		length:     int64(len(pcm)),
	}, nil
}

// readAll reads all PCM bytes from an io.Reader.
func readAll(r io.Reader) ([]byte, error) {
	buf := make([]byte, 8192)
	var all []byte
	for {
		n, readErr := r.Read(buf)
		if n > 0 {
			all = append(all, buf[:n]...)
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, fmt.Errorf("mp3: decode: %w", readErr)
		}
	}
	return all, nil
}

// DecodeWithSampleRate reads an MP3 file from src and resamples the
// output to the given sample rate.
func DecodeWithSampleRate(sampleRate int, src io.Reader) (*Stream, error) {
	s, err := Decode(src)
	if err != nil {
		return nil, err
	}
	if s.sampleRate == sampleRate {
		return s, nil
	}
	resampled := resample(s.raw, s.sampleRate, sampleRate)
	return &Stream{
		data:       bytes.NewReader(resampled),
		raw:        resampled,
		sampleRate: sampleRate,
		length:     int64(len(resampled)),
	}, nil
}

// resample performs linear interpolation resampling from srcRate to dstRate.
// Input and output are stereo 16-bit signed LE PCM.
func resample(data []byte, srcRate, dstRate int) []byte {
	srcFrames := len(data) / 4
	if srcFrames < 2 {
		return data
	}

	ratio := float64(srcRate) / float64(dstRate)
	dstFrames := int(float64(srcFrames) / ratio)
	out := make([]byte, dstFrames*4)

	for i := 0; i < dstFrames; i++ {
		srcPos := float64(i) * ratio
		srcIdx := int(srcPos)
		frac := srcPos - float64(srcIdx)

		if srcIdx >= srcFrames-1 {
			srcIdx = srcFrames - 2
			frac = 1.0
		}
		if srcIdx < 0 {
			srcIdx = 0
			frac = 0
		}

		l0 := int16(binary.LittleEndian.Uint16(data[srcIdx*4:]))
		l1 := int16(binary.LittleEndian.Uint16(data[(srcIdx+1)*4:]))
		left := int16(float64(l0) + frac*(float64(l1)-float64(l0)))

		r0 := int16(binary.LittleEndian.Uint16(data[srcIdx*4+2:]))
		r1 := int16(binary.LittleEndian.Uint16(data[(srcIdx+1)*4+2:]))
		right := int16(float64(r0) + frac*(float64(r1)-float64(r0)))

		binary.LittleEndian.PutUint16(out[i*4:], uint16(left))
		binary.LittleEndian.PutUint16(out[i*4+2:], uint16(right))
	}
	return out
}
