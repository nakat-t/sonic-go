package cgosonic

/*
#cgo CFLAGS: -Wall -Wno-unused-function -g -ansi -fPIC -pthread -I${SRCDIR}
#include <stdlib.h>
#include "sonic.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

const (
	MIN_VOLUME        = float32(C.SONIC_MIN_VOLUME)
	MAX_VOLUME        = float32(C.SONIC_MAX_VOLUME)
	MIN_SPEED         = float32(C.SONIC_MIN_SPEED)
	MAX_SPEED         = float32(C.SONIC_MAX_SPEED)
	MIN_PITCH_SETTING = float32(C.SONIC_MIN_PITCH_SETTING)
	MAX_PITCH_SETTING = float32(C.SONIC_MAX_PITCH_SETTING)
	MIN_RATE          = float32(C.SONIC_MIN_RATE)
	MAX_RATE          = float32(C.SONIC_MAX_RATE)
	MIN_SAMPLE_RATE   = int(C.SONIC_MIN_SAMPLE_RATE)
	MAX_SAMPLE_RATE   = int(C.SONIC_MAX_SAMPLE_RATE)
	MIN_CHANNELS      = int(C.SONIC_MIN_CHANNELS)
	MAX_CHANNELS      = int(C.SONIC_MAX_CHANNELS)
)

// Stream represents a SONIC audio stream
type Stream struct {
	stream C.sonicStream
}

// CreateStream creates a new sonic stream
func CreateStream(sampleRate int, numChannels int) (*Stream, error) {
	stream := C.sonicCreateStream(C.int(sampleRate), C.int(numChannels))
	if stream == nil {
		return nil, errors.New("failed to create cgosonic.Stream")
	}
	return &Stream{stream: stream}, nil
}

// DestroyStream destroys the sonic stream
func (s *Stream) DestroyStream() {
	if s.stream != nil {
		C.sonicDestroyStream(s.stream)
		s.stream = nil
	}
}

// The following symbols are not implemented yet.
// void sonicSetUserData(sonicStream stream, void *userData);
// void *sonicGetUserData(sonicStream stream);

// WriteFloatToStream writes float samples to the stream
func (s *Stream) WriteFloatToStream(samples []float32, numSamples int) int {
	return int(C.sonicWriteFloatToStream(s.stream, (*C.float)(unsafe.Pointer(&samples[0])), C.int(numSamples)))
}

// WriteShortToStream writes short samples to the stream
func (s *Stream) WriteShortToStream(samples []int16, numSamples int) int {
	return int(C.sonicWriteShortToStream(s.stream, (*C.short)(unsafe.Pointer(&samples[0])), C.int(numSamples)))
}

// The following symbol is not implemented yet.
// int sonicWriteUnsignedCharToStream(sonicStream stream, const unsigned char* samples, int numSamples);

// ReadFloatFromStream reads float samples from the stream
func (s *Stream) ReadFloatFromStream(samples []float32, maxSamples int) int {
	return int(C.sonicReadFloatFromStream(s.stream, (*C.float)(unsafe.Pointer(&samples[0])), C.int(maxSamples)))
}

// ReadShortFromStream reads short samples from the stream
func (s *Stream) ReadShortFromStream(samples []int16, maxSamples int) int {
	return int(C.sonicReadShortFromStream(s.stream, (*C.short)(unsafe.Pointer(&samples[0])), C.int(maxSamples)))
}

// The following symbol is not implemented yet.
// int sonicReadUnsignedCharFromStream(sonicStream stream, unsigned char* samples, int maxSamples);

// FlushStream flushes the stream
func (s *Stream) FlushStream() int {
	return int(C.sonicFlushStream(s.stream))
}

// SamplesAvailable returns the number of samples in the output buffer
func (s *Stream) SamplesAvailable() int {
	return int(C.sonicSamplesAvailable(s.stream))
}

// GetSpeed gets the speed of the stream
func (s *Stream) GetSpeed() float32 {
	return float32(C.sonicGetSpeed(s.stream))
}

// SetSpeed sets the speed of the stream
func (s *Stream) SetSpeed(speed float32) {
	C.sonicSetSpeed(s.stream, C.float(speed))
}

// GetPitch gets the pitch of the stream
func (s *Stream) GetPitch() float32 {
	return float32(C.sonicGetPitch(s.stream))
}

// SetPitch sets the pitch of the stream
func (s *Stream) SetPitch(pitch float32) {
	C.sonicSetPitch(s.stream, C.float(pitch))
}

// GetRate gets the rate of the stream
func (s *Stream) GetRate() float32 {
	return float32(C.sonicGetRate(s.stream))
}

// SetRate sets the rate of the stream
func (s *Stream) SetRate(rate float32) {
	C.sonicSetRate(s.stream, C.float(rate))
}

// GetVolume gets the volume of the stream
func (s *Stream) GetVolume() float32 {
	return float32(C.sonicGetVolume(s.stream))
}

// SetVolume sets the volume of the stream
func (s *Stream) SetVolume(volume float32) {
	C.sonicSetVolume(s.stream, C.float(volume))
}

// The following symbols are not implemented yet.
// int sonicGetChordPitch(sonicStream stream);
// void sonicSetChordPitch(sonicStream stream, int useChordPitch);

// GetQuality gets the quality setting.
func (s *Stream) GetQuality() int {
	return int(C.sonicGetQuality(s.stream))
}

// SetQuality sets the "quality".  Default 0 is virtually as good as 1, but very much faster.
func (s *Stream) SetQuality(quality int) {
	C.sonicSetQuality(s.stream, C.int(quality))
}

// GetSampleRate gets the sample rate of the stream
func (s *Stream) GetSampleRate() int {
	return int(C.sonicGetSampleRate(s.stream))
}

// SetSampleRate sets the sample rate of the stream
func (s *Stream) SetSampleRate(sampleRate int) {
	C.sonicSetSampleRate(s.stream, C.int(sampleRate))
}

// GetNumChannels gets the number of channels in the stream
func (s *Stream) GetNumChannels() int {
	return int(C.sonicGetNumChannels(s.stream))
}

// SetNumChannels sets the number of channels in the stream
func (s *Stream) SetNumChannels(numChannels int) {
	C.sonicSetNumChannels(s.stream, C.int(numChannels))
}

// ChangeFloatSpeed is a non-stream-oriented interface to change the speed of float audio samples
func ChangeFloatSpeed(samples []float32, numSamples int, speed, pitch, rate, volume float32, sampleRate, numChannels int) int {
	return int(C.sonicChangeFloatSpeed((*C.float)(unsafe.Pointer(&samples[0])), C.int(numSamples),
		C.float(speed), C.float(pitch), C.float(rate), C.float(volume),
		0, C.int(sampleRate), C.int(numChannels)))
}

// ChangeShortSpeed is a non-stream-oriented interface to change the speed of short audio samples
func ChangeShortSpeed(samples []int16, numSamples int, speed, pitch, rate, volume float32, sampleRate, numChannels int) int {
	return int(C.sonicChangeShortSpeed((*C.short)(unsafe.Pointer(&samples[0])), C.int(numSamples),
		C.float(speed), C.float(pitch), C.float(rate), C.float(volume),
		0, C.int(sampleRate), C.int(numChannels)))
}

// The following symbols are not implemented yet (SONIC_SPECTROGRAM related features).
// void sonicComputeSpectrogram(sonicStream stream);
// sonicSpectrogram sonicGetSpectrogram(sonicStream stream);
// sonicSpectrogram sonicCreateSpectrogram(int sampleRate);
// void sonicDestroySpectrogram(sonicSpectrogram spectrogram);
// sonicBitmap sonicConvertSpectrogramToBitmap(sonicSpectrogram spectrogram, int numRows, int numCols);
// void sonicDestroyBitmap(sonicBitmap bitmap);
// int sonicWritePGM(sonicBitmap bitmap, char* fileName);
// void sonicAddPitchPeriodToSpectrogram(sonicSpectrogram spectrogram, short* samples, int numSamples, int numChannels);
