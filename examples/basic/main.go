package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/nakat-t/sonic-go"
)

// Basic examples of using sonic

func main() {
	const sampleRate = 48000
	const bitsPerSample = 16
	const numChannels = 1
	const freq = 800
	const msec = 1000
	const amp = 16383 // 50% of int16 max value

	// Generate a beep sound
	src := GenerateBeep(sampleRate, freq, msec, amp)

	// Save source beep sound to a WAV file
	srcFile, _ := os.Create("src.wav")
	WriteWavHeader(srcFile, sampleRate, bitsPerSample, numChannels, src.Len())
	io.Copy(srcFile, src)
	srcFile.Close()

	out := bytes.NewBuffer(nil)

	// Re-generate the beep sound
	src = GenerateBeep(sampleRate, freq, msec, amp)

	// Create a Sonic transformer
	transformer, err := sonic.NewTransformer(out, sampleRate, sonic.AudioFormatPCM,
		sonic.WithVolume(0.2),
		sonic.WithSpeed(2.5),
		sonic.WithPitch(1.5),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	io.Copy(transformer, src)
	transformer.Flush()

	outFile, _ := os.Create("out.wav")
	WriteWavHeader(outFile, sampleRate, bitsPerSample, numChannels, out.Len())
	io.Copy(outFile, out)
	outFile.Close()
}

// GenerateBeep generates a sine wave beep sound
func GenerateBeep(sampleRate int, freq int, msec int, amp int) *bytes.Buffer {
	numSamples := sampleRate * msec / 1000
	buf := bytes.NewBuffer(nil)
	buf.Grow(numSamples * 2) // 2 bytes per sample for int16

	for i := 0; i < numSamples; i++ {
		// Time t (in seconds)
		t := float64(i) / float64(sampleRate)

		// Calculate sine waves and apply amplitude
		sample := int16(math.Round(float64(amp) * math.Sin(2.0*math.Pi*float64(freq)*t)))
		binary.Write(buf, binary.LittleEndian, sample)
	}

	return buf
}

func WriteWavHeader(w io.Writer, sampleRate int, bitsPerSample int, numChannels int, numDataBytes int) error {
	// WAV header size is 44 bytes
	header := make([]byte, 44)

	// RIFF header
	copy(header[0:4], []byte("RIFF"))
	binary.LittleEndian.PutUint32(header[4:8], uint32(numDataBytes+36))
	copy(header[8:12], []byte("WAVE"))

	// fmt subchunk
	copy(header[12:16], []byte("fmt "))
	binary.LittleEndian.PutUint32(header[16:20], 16) // Subchunk size for PCM
	binary.LittleEndian.PutUint16(header[20:22], 1)  // Audio format (PCM)
	binary.LittleEndian.PutUint16(header[22:24], uint16(numChannels))
	binary.LittleEndian.PutUint32(header[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(header[28:32], uint32(sampleRate*numChannels*bitsPerSample/8))
	binary.LittleEndian.PutUint16(header[32:34], uint16(numChannels*bitsPerSample/8))
	binary.LittleEndian.PutUint16(header[34:36], uint16(bitsPerSample))

	// data subchunk
	copy(header[36:40], []byte("data"))
	binary.LittleEndian.PutUint32(header[40:44], uint32(numDataBytes))

	_, err := w.Write(header)
	return err
}
