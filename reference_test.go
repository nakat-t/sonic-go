package sonic

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/nakat-t/sonic-go/internal/cgosonic"
)

const (
	// Path to the original audio file
	originalWavPath = "./test/testdata/original/common_voice_en_1dcef00e46910f33.wav"
	// Directory containing reference audio files
	referenceWavDir = "./test/testdata/reference/"
)

// TestReferenceVolume tests that volume modification matches the reference implementation
func TestReferenceVolume(t *testing.T) {
	volumeValues := []string{"0.01", "0.5", "1.0", "2.0", "100.0"}

	for _, volumeStr := range volumeValues {
		volume, err := strconv.ParseFloat(volumeStr, 32)
		if err != nil {
			t.Fatalf("Failed to parse: %v", err)
		}
		t.Run("Volume_"+volumeStr, func(t *testing.T) {
			testProcessedAudioMatchesReference(t, float32(volume), 1.0, 1.0, 0, "volume_"+volumeStr+".wav")
		})
	}
}

// TestReferenceSpeed tests that speed modification matches the reference implementation
func TestReferenceSpeed(t *testing.T) {
	speedValues := []string{"0.05", "0.5", "1.0", "2.0", "20.0"}

	for _, speedStr := range speedValues {
		speed, err := strconv.ParseFloat(speedStr, 32)
		if err != nil {
			t.Fatalf("Failed to parse: %v", err)
		}
		t.Run("Speed_"+speedStr, func(t *testing.T) {
			testProcessedAudioMatchesReference(t, 1.0, float32(speed), 1.0, 0, "speed_"+speedStr+".wav")
		})
	}
}

// TestReferencePitch tests that pitch modification matches the reference implementation
func TestReferencePitch(t *testing.T) {
	pitchValues := []string{"0.05", "0.5", "1.0", "2.0", "20.0"}

	for _, pitchStr := range pitchValues {
		pitch, err := strconv.ParseFloat(pitchStr, 32)
		if err != nil {
			t.Fatalf("Failed to parse: %v", err)
		}
		t.Run("Pitch_"+pitchStr, func(t *testing.T) {
			testProcessedAudioMatchesReference(t, 1.0, 1.0, float32(pitch), 0, "pitch_"+pitchStr+".wav")
		})
	}
}

// TestReferenceQuality tests that quality setting matches the reference implementation
func TestReferenceQuality(t *testing.T) {
	t.Run("Quality_On", func(t *testing.T) {
		testProcessedAudioMatchesReference(t, 1.0, 1.0, 1.0, 1, "quality_on.wav")
	})
}

// testProcessedAudioMatchesReference verifies that audio processed with specified parameters matches the reference audio
func testProcessedAudioMatchesReference(t *testing.T, volume, speed, pitch float32, quality int, referenceFileName string) {
	t.Helper()
	t.Logf("Testing volume=%v, speed=%v, pitch=%v, quality=%v, file=%v", volume, speed, pitch, quality, referenceFileName)

	const BUFFER_SIZE = 4096

	// Load the original audio file
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	wfIn, sampleRate, numChannels, err := cgosonic.OpenInputWaveFile(filepath.Join(cwd, originalWavPath))
	if err != nil {
		t.Fatalf("Failed to open original audio file: %v", err)
	}
	wfIn.CloseWaveFile()

	fileIn, err := os.Open(filepath.Join(cwd, originalWavPath))
	if err != nil {
		t.Fatalf("Failed to open original audio file: %v", err)
	}
	in := bytes.NewBuffer(nil)
	_, err = io.Copy(in, fileIn)
	if err != nil {
		t.Fatalf("Failed to read original audio file: %v", err)
	}
	fileIn.Close()
	in.Next(44) // Skip the WAV header

	opts := []Option{
		WithSpeed(speed),
		WithPitch(pitch),
		WithVolume(volume),
	}
	if quality != 0 {
		opts = append(opts, WithQuality())
	}
	if numChannels != 1 {
		opts = append(opts, WithChannels(numChannels))
	}

	out := bytes.NewBuffer(nil)

	// Create a Sonic instance
	transformer, err := NewTransformer(out, sampleRate, AudioFormatPCM, opts...)
	if err != nil {
		t.Fatalf("Failed to create Sonic instance: %v", err)
	}

	_, err = io.Copy(transformer, in)
	if err != nil {
		t.Fatalf("Failed to copy data to transformer: %v", err)
	}

	transformer.Flush()

	processedSamples := make([]int16, out.Len()/2)
	if err := binary.Read(out, binary.LittleEndian, processedSamples); err != nil {
		t.Fatalf("Failed to read processed samples: %v", err)
	}

	// For Debug: Output processed wave file to 'test/testdata/processed/sonic/'
	if os.Getenv("CGOSONIC_TEST_DEBUG") != "" {
		os.MkdirAll(filepath.Join(cwd, "./test/testdata/processed/sonic/"), 0755)

		processedWavPath := filepath.Join(cwd, "./test/testdata/processed/sonic/", referenceFileName)
		wfOut, err := cgosonic.OpenOutputWaveFile(processedWavPath, sampleRate, numChannels)
		if err != nil {
			t.Fatalf("Failed to open output wave file: %v", err)
		}
		defer wfOut.CloseWaveFile()

		okWritten := wfOut.WriteToWaveFile(processedSamples, len(processedSamples))
		if okWritten == 0 {
			t.Errorf("Failed to write all samples to output wave file")
		}
	}

	numReferenceSamples, err := readNumSamplesFromWavFile(filepath.Join(cwd, referenceWavDir, referenceFileName))
	if err != nil {
		t.Logf("If you don't have any reference audio yet, run the following command: ./scripts/gen-testdata.sh")
		t.Fatalf("Failed to read reference sample count: %v", err)
	}

	// Load reference audio file
	referenceFilePath := filepath.Join(referenceWavDir, referenceFileName)
	refWf, refSampleRate, refNumChannels, err := cgosonic.OpenInputWaveFile(referenceFilePath)
	if err != nil {
		t.Fatalf("Failed to open reference audio file: %v", err)
	}
	defer refWf.CloseWaveFile()

	// Verify sample rate and channel count of reference audio
	if refSampleRate != sampleRate {
		t.Errorf("Reference audio has different sample rate: %d != %d", refSampleRate, sampleRate)
	}
	if refNumChannels != numChannels {
		t.Errorf("Reference audio has different channel count: %d != %d", refNumChannels, numChannels)
	}

	// Read reference audio samples
	referenceBufferSize := numReferenceSamples + BUFFER_SIZE
	referenceBuffer := make([]int16, referenceBufferSize)
	numReferenceSamplesRead := 0
	for {
		// Read samples from the WAVE file
		samplesRead := refWf.ReadFromWaveFile(referenceBuffer[numReferenceSamplesRead:], referenceBufferSize-numReferenceSamplesRead)
		if samplesRead <= 0 {
			break // No more samples to read
		}
		numReferenceSamplesRead += samplesRead
	}
	if numReferenceSamplesRead <= 0 {
		t.Fatalf("Failed to read reference audio samples: %d", numReferenceSamplesRead)
	}
	if numReferenceSamplesRead != numReferenceSamples {
		t.Fatalf("Reference sample count mismatch: expected %d, got %d", numReferenceSamples, numReferenceSamplesRead)
	}

	// Trim reference buffer to actual size read
	referenceBuffer = referenceBuffer[:numReferenceSamples]

	// Compare sample counts
	samplesAllowedDiffPercent := 1.0 // Allowable difference in sample count (1% of the smaller buffer)
	samplesDiffPercent := float64(abs(len(processedSamples)-len(referenceBuffer))) / float64(min(len(processedSamples), len(referenceBuffer))) * 100.0
	if samplesDiffPercent > samplesAllowedDiffPercent {
		t.Errorf("Processed sample count differs from reference sample count: %.2f%% > %.2f%%",
			samplesDiffPercent, samplesAllowedDiffPercent)
		t.Logf("Processed samples: %d, Reference samples: %d", len(processedSamples), len(referenceBuffer))
	}

	// Compare sample content
	// Due to potential minor differences in WAV encoding/decoding precision,
	// check if the difference is within tolerance rather than requiring exact match
	maxDiff := 0
	maxDiffIndex := -1
	maxSamplesToCompare := min(len(processedSamples), len(referenceBuffer))

	const toleranceDiff = 5 // Tolerance threshold (absolute value)
	differenceCount := 0

	for i := range maxSamplesToCompare {
		diff := abs(int(processedSamples[i]) - int(referenceBuffer[i]))
		if diff > maxDiff {
			maxDiff = diff
			maxDiffIndex = i
		}

		if diff > toleranceDiff {
			differenceCount++
		}
	}

	// Test passes if percentage of samples outside tolerance is less than 1%
	maxAllowedDiffPercent := 1.0
	diffPercent := float64(differenceCount) / float64(maxSamplesToCompare) * 100.0

	if diffPercent > maxAllowedDiffPercent {
		t.Errorf("Difference between processed and reference samples exceeds tolerance: %.2f%% > %.2f%%",
			diffPercent, maxAllowedDiffPercent)
		t.Logf("Maximum difference: %d (index %d, processed=%d, reference=%d)",
			maxDiff, maxDiffIndex, processedSamples[maxDiffIndex], referenceBuffer[maxDiffIndex])
	} else {
		t.Logf("Sample comparison result: difference %.2f%%, maximum difference %d", diffPercent, maxDiff)
	}
}

// Helper functions

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// readNumSamplesFromWavFile reads the number of samples from a WAV file
func readNumSamplesFromWavFile(filePath string) (int, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	fileSize := stat.Size()
	dataSize := fileSize - 44     // Subtract header size (44 bytes for WAV)
	return int(dataSize / 2), nil // Assuming 16-bit PCM (2 bytes per sample)
}
