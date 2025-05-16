package sonic

import (
	"testing"

	"github.com/nakat-t/sonic-go/internal/cgosonic"
)

// Note: These tests assume that the Transformer struct is defined elsewhere in the 'sonic' package
// and has at least the following fields with the specified types:
// type Transformer struct {
//     numChannels int
//     volume      *float32
//     speed       *float32
//     pitch       *float32
//     rate        *float32
//     quality     *int
//     // ... other fields
// }

func TestWithChannels(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"within range (2)", 2, 2},
		{"within range (1)", 1, 1},
		{"within range (32)", 32, 32},
		{"below min", cgosonic.MIN_CHANNELS - 1, cgosonic.MIN_CHANNELS},
		{"at min", cgosonic.MIN_CHANNELS, cgosonic.MIN_CHANNELS},
		{"above max", cgosonic.MAX_CHANNELS + 1, cgosonic.MAX_CHANNELS},
		{"at max", cgosonic.MAX_CHANNELS, cgosonic.MAX_CHANNELS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{} // Assuming Transformer struct exists
			opt := WithChannels(tt.input)
			err := opt(tr)
			if err != nil {
				t.Fatalf("WithChannels(%d) returned an error: %v", tt.input, err)
			}
			if tr.numChannels != tt.expected {
				t.Errorf("WithChannels(%d) set numChannels to %d; want %d", tt.input, tr.numChannels, tt.expected)
			}
		})
	}
}

func TestWithVolume(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{"within range (1.5)", 1.5, 1.5},
		{"below min", cgosonic.MIN_VOLUME - 0.01, cgosonic.MIN_VOLUME},
		{"at min", cgosonic.MIN_VOLUME, cgosonic.MIN_VOLUME},
		{"above max", cgosonic.MAX_VOLUME + 0.1, cgosonic.MAX_VOLUME},
		{"at max", cgosonic.MAX_VOLUME, cgosonic.MAX_VOLUME},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{}
			opt := WithVolume(tt.input)
			err := opt(tr)
			if err != nil {
				t.Fatalf("WithVolume(%f) returned an error: %v", tt.input, err)
			}
			if tr.volume == nil {
				t.Fatalf("WithVolume(%f) did not set volume, field is nil", tt.input)
			}
			if *tr.volume != tt.expected {
				t.Errorf("WithVolume(%f) set volume to %f; want %f", tt.input, *tr.volume, tt.expected)
			}
		})
	}
}

func TestWithSpeed(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{"within range (2.0)", 2.0, 2.0},
		{"below min", cgosonic.MIN_SPEED - 0.01, cgosonic.MIN_SPEED},
		{"at min", cgosonic.MIN_SPEED, cgosonic.MIN_SPEED},
		{"above max", cgosonic.MAX_SPEED + 0.1, cgosonic.MAX_SPEED},
		{"at max", cgosonic.MAX_SPEED, cgosonic.MAX_SPEED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{}
			opt := WithSpeed(tt.input)
			err := opt(tr)
			if err != nil {
				t.Fatalf("WithSpeed(%f) returned an error: %v", tt.input, err)
			}
			if tr.speed == nil {
				t.Fatalf("WithSpeed(%f) did not set speed, field is nil", tt.input)
			}
			if *tr.speed != tt.expected {
				t.Errorf("WithSpeed(%f) set speed to %f; want %f", tt.input, *tr.speed, tt.expected)
			}
		})
	}
}

func TestWithPitch(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{"within range (1.3)", 1.3, 1.3},
		{"below min", cgosonic.MIN_PITCH_SETTING - 0.01, cgosonic.MIN_PITCH_SETTING},
		{"at min", cgosonic.MIN_PITCH_SETTING, cgosonic.MIN_PITCH_SETTING},
		{"above max", cgosonic.MAX_PITCH_SETTING + 0.1, cgosonic.MAX_PITCH_SETTING},
		{"at max", cgosonic.MAX_PITCH_SETTING, cgosonic.MAX_PITCH_SETTING},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{}
			opt := WithPitch(tt.input)
			err := opt(tr)
			if err != nil {
				t.Fatalf("WithPitch(%f) returned an error: %v", tt.input, err)
			}
			if tr.pitch == nil {
				t.Fatalf("WithPitch(%f) did not set pitch, field is nil", tt.input)
			}
			if *tr.pitch != tt.expected {
				t.Errorf("WithPitch(%f) set pitch to %f; want %f", tt.input, *tr.pitch, tt.expected)
			}
		})
	}
}

func TestWithRate(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		expected float32
	}{
		{"within range (2.0)", 2.0, 2.0},
		{"below min", cgosonic.MIN_RATE - 0.01, cgosonic.MIN_RATE},
		{"at min", cgosonic.MIN_RATE, cgosonic.MIN_RATE},
		{"above max", cgosonic.MAX_RATE + 0.1, cgosonic.MAX_RATE},
		{"at max", cgosonic.MAX_RATE, cgosonic.MAX_RATE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Transformer{}
			opt := WithRate(tt.input)
			err := opt(tr)
			if err != nil {
				t.Fatalf("WithRate(%f) returned an error: %v", tt.input, err)
			}
			if tr.rate == nil {
				t.Fatalf("WithRate(%f) did not set rate, field is nil", tt.input)
			}
			if *tr.rate != tt.expected {
				t.Errorf("WithRate(%f) set rate to %f; want %f", tt.input, *tr.rate, tt.expected)
			}
		})
	}
}

func TestWithQuality(t *testing.T) {
	tr := &Transformer{}
	opt := WithQuality()
	err := opt(tr)
	if err != nil {
		t.Fatalf("WithQuality() returned an error: %v", err)
	}
	if tr.quality == nil {
		t.Fatal("WithQuality() did not set quality, field is nil")
	}
	if *tr.quality != 1 {
		t.Errorf("WithQuality() set quality to %d; want 1", *tr.quality)
	}
}
