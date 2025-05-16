// Package sonic is Go bindings for libsonic. Sonic is a algorithm for changing speech speed, uniquely optimized for over 2X acceleration.
package sonic

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"runtime"
	"unsafe"

	"github.com/nakat-t/sonic-go/internal/cgosonic"
)

var (
	// ErrInvalid is returned when an invalid value is provided.
	ErrInvalid = errors.New("invalid value")

	// ErrWrite is returned when writing to the writer fails.
	ErrWrite = errors.New("failed to write to writer")

	// ErrSonicCLib is returned when Sonic C library fails.
	ErrSonicCLib = errors.New("sonic C library error")

	// ErrInternal is returned when an internal error occurs.
	ErrInternal = errors.New("internal error")
)

// Constants for audio formats
const (
	FormatInt16   = 1 // PCM 16-bit signed integer
	FormatFloat32 = 3 // PCM 32-bit float
)

const (
	streamBufferSize = 4096 // Buffer size for cgosonic.Stream
)

// Transformer is a struct that transforms audio data using the Sonic library.
type Transformer struct {
	w           io.Writer
	sampleRate  int
	numChannels int
	format      int
	volume      *float32
	speed       *float32
	pitch       *float32
	rate        *float32
	quality     *int

	stream       *cgosonic.Stream
	streamBuffer []byte
}

// NewTransformer creates a new Transformer instance.
func NewTransformer(w io.Writer, sampleRate int, format int, opts ...Option) (*Transformer, error) {
	if w == nil {
		return nil, fmt.Errorf("%w: writer is nil", ErrInvalid)
	}
	if sampleRate < cgosonic.MIN_SAMPLE_RATE || cgosonic.MAX_SAMPLE_RATE < sampleRate {
		return nil, fmt.Errorf("%w: sampleRate %d is out of range [%d, %d]", ErrInvalid, sampleRate, cgosonic.MIN_SAMPLE_RATE, cgosonic.MAX_SAMPLE_RATE)
	}
	if format != FormatInt16 && format != FormatFloat32 {
		return nil, fmt.Errorf("%w: format %d is not supported", ErrInvalid, format)
	}

	t := &Transformer{
		w:            w,
		sampleRate:   sampleRate,
		numChannels:  1,
		format:       format,
		volume:       nil,
		speed:        nil,
		pitch:        nil,
		rate:         nil,
		quality:      nil,
		stream:       nil,
		streamBuffer: nil,
	}
	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}

	stream, err := cgosonic.CreateStream(t.sampleRate, t.numChannels)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create sonic stream", ErrSonicCLib)
	}
	t.stream = stream

	switch t.format {
	case FormatInt16:
		t.streamBuffer = make([]byte, streamBufferSize*2) // 2 bytes per sample for int16
	case FormatFloat32:
		t.streamBuffer = make([]byte, streamBufferSize*4) // 4 bytes per sample for float32
	default:
		return nil, fmt.Errorf("%w: format is broken: %d", ErrInternal, t.format)
	}

	if t.volume != nil {
		stream.SetVolume(*t.volume)
	}
	if t.speed != nil {
		stream.SetSpeed(*t.speed)
	}
	if t.pitch != nil {
		stream.SetPitch(*t.pitch)
	}
	if t.rate != nil {
		stream.SetRate(*t.rate)
	}
	if t.quality != nil {
		stream.SetQuality(*t.quality)
	}

	runtime.SetFinalizer(t, func(t *Transformer) {
		if t != nil {
			t.Close()
		}
	})

	return t, nil
}

// Write writes the data to the transformer.
func (t *Transformer) Write(p []byte) (int, error) {
	switch t.format {
	case FormatInt16:
		return t.writeInt16(p)
	case FormatFloat32:
		return t.writeFloat32(p)
	default:
		return 0, fmt.Errorf("%w: format is broken: %d", ErrInternal, t.format)
	}
}

// Flush flushes the transformer.
func (t *Transformer) Flush() error {
	switch t.format {
	case FormatInt16:
		return t.flushInt16()
	case FormatFloat32:
		return t.flushFloat32()
	default:
		return fmt.Errorf("%w: format is broken: %d", ErrInternal, t.format)
	}
}

// Close closes the transformer and releases resources.
func (t *Transformer) Close() error {
	if t.stream != nil {
		t.stream.DestroyStream()
		t.stream = nil
	}
	if t.streamBuffer != nil {
		t.streamBuffer = nil
	}
	return nil
}

// writeInt16 writes int16 data to the transformer.
func (t *Transformer) writeInt16(p []byte) (int, error) {
	const sampleSize = 2                                         // 2 bytes per sample for int16
	const streamBufferSampleSize = streamBufferSize / sampleSize // Number of samples in the stream buffer

	if len(p)%sampleSize != 0 {
		return 0, fmt.Errorf("%w: 'p' must be a multiple of the int16 type size", ErrInvalid)
	}
	samples := t.unsafeBytesAsInt16Slice(p)
	if len(samples) == 0 {
		return 0, nil
	}

	numWrittenBytes := 0

	for {
		size := min(len(samples), streamBufferSampleSize)
		if size <= 0 {
			break
		}
		okInt := t.stream.WriteShortToStream(samples[:size], size/t.numChannels)
		if okInt == 0 {
			return numWrittenBytes, fmt.Errorf("%w: failed to write samples to stream", ErrSonicCLib)
		}
		numWrittenBytes += size * sampleSize

		buf := t.unsafeBytesAsInt16Slice(t.streamBuffer)
		for {
			nRead := t.stream.ReadShortFromStream(buf, len(buf)/t.numChannels)
			if nRead <= 0 {
				break
			}
			if err := binary.Write(t.w, binary.LittleEndian, buf[:nRead]); err != nil {
				return numWrittenBytes, fmt.Errorf("%w: failed to write samples: %w", ErrWrite, err)
			}
		}

		samples = samples[size:]
	}

	return numWrittenBytes, nil
}

// writeFloat32 writes float32 data to the transformer.
func (t *Transformer) writeFloat32(p []byte) (int, error) {
	const sampleSize = 4                                         // 4 bytes per sample for float32
	const streamBufferSampleSize = streamBufferSize / sampleSize // Number of samples in the stream buffer

	if len(p)%sampleSize != 0 {
		return 0, fmt.Errorf("%w: 'p' must be a multiple of the float32 type size", ErrInvalid)
	}
	samples := t.unsafeBytesAsFloat32Slice(p)
	if len(samples) == 0 {
		return 0, nil
	}

	numWrittenBytes := 0

	for {
		size := min(len(samples), streamBufferSampleSize)
		if size <= 0 {
			break
		}
		okInt := t.stream.WriteFloatToStream(samples[:size], size/t.numChannels)
		if okInt == 0 {
			return numWrittenBytes, fmt.Errorf("%w: failed to write samples to stream", ErrSonicCLib)
		}
		numWrittenBytes += size * sampleSize

		buf := t.unsafeBytesAsFloat32Slice(t.streamBuffer)
		for {
			nRead := t.stream.ReadFloatFromStream(buf, len(buf)/t.numChannels)
			if nRead <= 0 {
				break
			}
			if err := binary.Write(t.w, binary.LittleEndian, buf[:nRead]); err != nil {
				return numWrittenBytes, fmt.Errorf("%w: failed to write samples: %w", ErrWrite, err)
			}
		}

		samples = samples[size:]
	}

	return numWrittenBytes, nil
}

func (t *Transformer) flushInt16() error {
	ret := t.stream.FlushStream()
	if ret == 0 {
		return fmt.Errorf("%w: failed to flush stream", ErrSonicCLib)
	}
	for t.stream.SamplesAvailable() > 0 {
		samples := make([]int16, t.stream.SamplesAvailable())
		n := t.stream.ReadShortFromStream(samples, len(samples))
		if n <= 0 {
			return fmt.Errorf("%w: failed to read samples from stream", ErrSonicCLib)
		}
		if err := binary.Write(t.w, binary.LittleEndian, samples[:n]); err != nil {
			return fmt.Errorf("%w: failed to write samples: %w", ErrWrite, err)
		}
	}
	return nil
}

func (t *Transformer) flushFloat32() error {
	ret := t.stream.FlushStream()
	if ret == 0 {
		return fmt.Errorf("%w: failed to flush stream", ErrSonicCLib)
	}
	for t.stream.SamplesAvailable() > 0 {
		samples := make([]float32, t.stream.SamplesAvailable())
		n := t.stream.ReadFloatFromStream(samples, len(samples))
		if n <= 0 {
			return fmt.Errorf("%w: failed to read samples from stream", ErrSonicCLib)
		}
		if err := binary.Write(t.w, binary.LittleEndian, samples[:n]); err != nil {
			return fmt.Errorf("%w: failed to write samples: %w", ErrWrite, err)
		}
	}
	return nil
}

func (t *Transformer) unsafeBytesAsInt16Slice(p []byte) []int16 {
	numSamples := len(p) / 2 // 2 bytes per sample for int16
	if numSamples == 0 {
		return nil
	}
	return (*[1 << 30]int16)(unsafe.Pointer(&p[0]))[:numSamples]
}

func (t *Transformer) unsafeBytesAsFloat32Slice(p []byte) []float32 {
	numSamples := len(p) / 4 // 4 bytes per sample for float32
	if numSamples == 0 {
		return nil
	}
	return (*[1 << 30]float32)(unsafe.Pointer(&p[0]))[:numSamples]
}
