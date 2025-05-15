package sonic

import (
	"cmp"

	"github.com/nakat-t/sonic-go/internal/cgosonic"
)

type Option func(*Transformer) error

// WithChannels sets the number of channels.
//
// You can specify a value between 1 and 32. Values outside this range are clamped.
// The default value is 1 (mono).
func WithChannels(channels int) Option {
	return func(t *Transformer) error {
		t.numChannels = clamp(channels, cgosonic.MIN_CHANNELS, cgosonic.MAX_CHANNELS)
		return nil
	}
}

// WithVolume sets the volume.
//
// This value scales the volume by a constant factor.
// You can specify a value between 0.01 and 100. Values outside this range are clamped.
// The default value is 1.0.
func WithVolume(volume float32) Option {
	return func(t *Transformer) error {
		val := clamp(volume, cgosonic.MIN_VOLUME, cgosonic.MAX_VOLUME)
		t.volume = &val
		return nil
	}
}

// WithSpeed sets the speed up factor.
//
// This value scales the speed. 2.0 means 2X faster.
// You can specify a value between 0.05 and 20. Values outside this range are clamped.
// The default value is 1.0.
func WithSpeed(speed float32) Option {
	return func(t *Transformer) error {
		val := clamp(speed, cgosonic.MIN_SPEED, cgosonic.MAX_SPEED)
		t.speed = &val
		return nil
	}
}

// WithPitch sets the pitch scaling factor.
//
// This value scales the pitch. 1.3 means 30% higher.
// You can specify a value between 0.05 and 20. Values outside this range are clamped.
// The default value is 1.0.
func WithPitch(pitch float32) Option {
	return func(t *Transformer) error {
		val := clamp(pitch, cgosonic.MIN_PITCH_SETTING, cgosonic.MAX_PITCH_SETTING)
		t.pitch = &val
		return nil
	}
}

// WithRate sets the playback rate.
//
// This value scales the playback rate. 2.0 means 2X faster, and 2X pitch.
// You can specify a value between 0.05 and 20. Values outside this range are clamped.
// The default value is 1.0.
func WithRate(rate float32) Option {
	return func(t *Transformer) error {
		val := clamp(rate, cgosonic.MIN_RATE, cgosonic.MAX_RATE)
		t.rate = &val
		return nil
	}
}

// WithQuality sets the quality.
//
// Setting the 'quality' flag disables speed-up heuristics. May increase quality.
// The default value is OFF (= Enable speed-up heuristics).
// Default OFF is virtually as good as ON (= Disable speed-up heuristics), but very much faster.
func WithQuality() Option {
	return func(t *Transformer) error {
		val := 1
		t.quality = &val
		return nil
	}
}

func clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
