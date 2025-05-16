package cgosonic

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// createDummyWav creates a minimal WAV file for testing.
// numSampleFrames is the number of sample frames.
// bitsPerSample is typically 16 for int16 samples.
func createDummyWav(t *testing.T, filename string, sampleRate int, numChannels int, numSampleFrames int, bitsPerSample int) {
	t.Helper()

	bytesPerSample := bitsPerSample / 8
	dataSize := uint32(numSampleFrames * numChannels * bytesPerSample)
	// RIFF Chunk Size = 4 (WAVE) + (8 (fmt header) + 16 (fmt data)) + (8 (data header) + dataSize)
	// = 4 + 24 + 8 + dataSize = 36 + dataSize
	riffChunkSize := uint32(36 + dataSize)

	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create dummy wav file %s: %v", filename, err)
	}
	defer file.Close()

	// RIFF Chunk Descriptor
	if _, err := file.WriteString("RIFF"); err != nil {
		t.Fatalf("Error writing RIFF: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, riffChunkSize); err != nil {
		t.Fatalf("Error writing riffChunkSize: %v", err)
	}
	if _, err := file.WriteString("WAVE"); err != nil {
		t.Fatalf("Error writing WAVE: %v", err)
	}

	// fmt Sub-chunk
	if _, err := file.WriteString("fmt "); err != nil {
		t.Fatalf("Error writing fmt : %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(16)); err != nil { // Subchunk1Size for PCM
		t.Fatalf("Error writing Subchunk1Size: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(1)); err != nil { // AudioFormat (1 for PCM)
		t.Fatalf("Error writing AudioFormat: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(numChannels)); err != nil {
		t.Fatalf("Error writing numChannels: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(sampleRate)); err != nil {
		t.Fatalf("Error writing sampleRate: %v", err)
	}
	byteRate := uint32(sampleRate * numChannels * bytesPerSample)
	if err := binary.Write(file, binary.LittleEndian, byteRate); err != nil {
		t.Fatalf("Error writing byteRate: %v", err)
	}
	blockAlign := uint16(numChannels * bytesPerSample)
	if err := binary.Write(file, binary.LittleEndian, blockAlign); err != nil {
		t.Fatalf("Error writing blockAlign: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(bitsPerSample)); err != nil {
		t.Fatalf("Error writing bitsPerSample: %v", err)
	}

	// data Sub-chunk
	if _, err := file.WriteString("data"); err != nil {
		t.Fatalf("Error writing data: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, dataSize); err != nil {
		t.Fatalf("Error writing dataSize: %v", err)
	}

	// Actual sample data (zeros)
	dummySampleData := make([]byte, dataSize)
	if _, err := file.Write(dummySampleData); err != nil {
		t.Fatalf("Failed to write dummy sample data: %v", err)
	}
}

func TestOpenInputWaveFile(t *testing.T) {
	tempDir := t.TempDir()
	validFileName := filepath.Join(tempDir, "test_input.wav")
	nonExistentFileName := filepath.Join(tempDir, "non_existent.wav")

	expectedSampleRate := 44100
	expectedNumChannels := 1
	createDummyWav(t, validFileName, expectedSampleRate, expectedNumChannels, 100, 16)

	t.Run("Success", func(t *testing.T) {
		wf, sr, nc, err := OpenInputWaveFile(validFileName)
		if err != nil {
			t.Fatalf("OpenInputWaveFile failed for valid file: %v", err)
		}
		if wf == nil {
			t.Fatal("OpenInputWaveFile returned nil WaveFile for valid file")
		}
		if wf.file == nil {
			t.Fatal("WaveFile.file is nil after successful open")
		}
		if sr != expectedSampleRate {
			t.Errorf("Expected sampleRate %d, got %d", expectedSampleRate, sr)
		}
		if nc != expectedNumChannels {
			t.Errorf("Expected numChannels %d, got %d", expectedNumChannels, nc)
		}

		if wf != nil {
			// Assuming 0 is success for C.closeWaveFile.
			// The actual success/error codes depend on the C implementation.
			res := wf.CloseWaveFile()
			if res != 0 {
				t.Logf("CloseWaveFile returned non-zero: %d. This might be an error or a status code from the C library.", res)
			}
		}
	})

	t.Run("ErrorNonExistentFile", func(t *testing.T) {
		wf, sr, nc, err := OpenInputWaveFile(nonExistentFileName)
		if err == nil {
			t.Fatal("OpenInputWaveFile succeeded for non-existent file, expected error")
		}
		if wf != nil {
			t.Error("OpenInputWaveFile returned non-nil WaveFile for non-existent file")
			if wf.file != nil {
				wf.CloseWaveFile() // Attempt to clean up if C layer somehow opened it
			}
		}
		if sr != 0 {
			t.Errorf("Expected sampleRate to be 0 on error, got %d", sr)
		}
		if nc != 0 {
			t.Errorf("Expected numChannels to be 0 on error, got %d", nc)
		}
	})
}

func TestOpenOutputWaveFile(t *testing.T) {
	tempDir := t.TempDir()
	validFileName := filepath.Join(tempDir, "test_output.wav")

	expectedSampleRate := 44100
	expectedNumChannels := 1
	createDummyWav(t, validFileName, expectedSampleRate, expectedNumChannels, 100, 16)

	t.Run("Success", func(t *testing.T) {
		wf, err := OpenOutputWaveFile(validFileName, expectedSampleRate, expectedNumChannels)
		if err != nil {
			t.Fatalf("OpenOutputWaveFile failed: %v", err)
		}
		if wf == nil {
			t.Fatal("OpenOutputWaveFile returned nil WaveFile")
		}
		if wf.file == nil {
			t.Fatal("WaveFile.file is nil after successful open")
		}
		defer os.Remove(validFileName) // Clean up the created file

		res := wf.CloseWaveFile()
		if res != 0 {
			t.Logf("CloseWaveFile returned non-zero: %d.", res)
		}
	})

	t.Run("ErrorInvalidPath", func(t *testing.T) {
		// Attempt to create a file in a non-existent directory
		invalidPath := filepath.Join(tempDir, "non_existent_dir", "test.wav")
		wf, err := OpenOutputWaveFile(invalidPath, 44100, 1)
		if err == nil {
			t.Fatal("OpenOutputWaveFile succeeded for invalid path, expected error")
			if wf != nil {
				if wf.file != nil {
					wf.CloseWaveFile()
				}
				os.Remove(invalidPath) // Attempt to clean up if created
			}
		}
		if wf != nil {
			t.Error("OpenOutputWaveFile returned non-nil WaveFile for invalid path")
		}
	})
}

func TestWaveFile_CloseWaveFile(t *testing.T) {
	tempDir := t.TempDir()
	fileName := filepath.Join(tempDir, "test_close.wav")

	expectedSampleRate := 44100
	expectedNumChannels := 1
	createDummyWav(t, fileName, expectedSampleRate, expectedNumChannels, 100, 16)

	wf, err := OpenOutputWaveFile(fileName, expectedSampleRate, expectedNumChannels)
	if err != nil {
		t.Fatalf("Setup: OpenOutputWaveFile failed: %v", err)
	}
	if wf == nil {
		t.Fatal("Setup: OpenOutputWaveFile returned nil WaveFile")
	}
	defer os.Remove(fileName)

	result := wf.CloseWaveFile()
	if result != 1 { // Assuming 0 is success for C.closeWaveFile
		t.Errorf("CloseWaveFile returned %d, expected 1 for success", result)
	}
	if wf.file != nil {
		t.Error("WaveFile.file is not nil after CloseWaveFile")
	}

	// Test closing an already closed file (wf.file is now nil)
	// The C function `closeWaveFile` will receive a NULL pointer.
	// Behavior depends on C library (could be no-op, error, or crash).
	resultAlreadyClosed := wf.CloseWaveFile()
	t.Logf("Closing an already closed file (nil C pointer) returned: %d. This behavior depends on the C library.", resultAlreadyClosed)
	if wf.file != nil {
		// This should not happen as Go wrapper sets wf.file to nil and doesn't reset it.
		t.Error("WaveFile.file is not nil after second CloseWaveFile call (should remain nil)")
	}
}

func TestWaveFile_ReadWriteToWaveFile(t *testing.T) {
	tempDir := t.TempDir()
	inputFileName := filepath.Join(tempDir, "rw_input.wav")
	outputFileName := filepath.Join(tempDir, "rw_output.wav")

	inputSampleRate := 8000
	inputNumChannels := 1
	inputNumSampleFrames := 256 // Number of sample frames in the input file (16-bit PCM)

	createDummyWav(t, inputFileName, inputSampleRate, inputNumChannels, inputNumSampleFrames, 16)
	defer os.Remove(inputFileName)

	outputSampleRate := 8000
	outputNumChannels := 1
	outputNumSampleFrames := 256 // Number of sample frames in the output file (16-bit PCM)
	createDummyWav(t, outputFileName, outputSampleRate, outputNumChannels, outputNumSampleFrames, 16)
	defer os.Remove(outputFileName)

	inputFile, sr, nc, err := OpenInputWaveFile(inputFileName)
	if err != nil {
		t.Fatalf("Failed to open input wave file '%s': %v", inputFileName, err)
	}
	if inputFile == nil {
		t.Fatal("OpenInputWaveFile returned nil for input")
	}
	defer inputFile.CloseWaveFile()

	if sr != inputSampleRate {
		t.Errorf("Input file sample rate mismatch: expected %d, got %d", inputSampleRate, sr)
	}
	if nc != inputNumChannels {
		t.Errorf("Input file num channels mismatch: expected %d, got %d", inputNumChannels, nc)
	}

	outputFile, err := OpenOutputWaveFile(outputFileName, sr, nc) // Use sr, nc from input
	if err != nil {
		t.Fatalf("Failed to open output wave file '%s': %v", outputFileName, err)
	}
	if outputFile == nil {
		t.Fatal("OpenOutputWaveFile returned nil for output")
	}
	defer outputFile.CloseWaveFile()

	// Read/write in chunks
	bufferFramesPerChunk := 128                       // Number of sample frames to process per chunk
	bufferShortsPerChunk := bufferFramesPerChunk * nc // Total int16 values in the buffer for one chunk
	buffer := make([]int16, bufferShortsPerChunk)

	totalShortsRead := 0
	expectedTotalShortsToRead := inputNumSampleFrames * nc

	for totalShortsRead < expectedTotalShortsToRead {
		shortsToReadInThisCall := bufferShortsPerChunk
		if remainingShorts := expectedTotalShortsToRead - totalShortsRead; remainingShorts < shortsToReadInThisCall {
			shortsToReadInThisCall = remainingShorts
		}

		// ReadFromWaveFile's maxSamples is total shorts. Returns shorts read.
		shortsRead := inputFile.ReadFromWaveFile(buffer, shortsToReadInThisCall)
		if shortsRead < 0 {
			t.Fatalf("ReadFromWaveFile returned negative value: %d", shortsRead)
		}
		if shortsRead == 0 {
			if totalShortsRead < expectedTotalShortsToRead {
				t.Errorf("ReadFromWaveFile returned 0 shorts prematurely. Read %d of %d expected.", totalShortsRead, expectedTotalShortsToRead)
			}
			break // EOF or error
		}
		totalShortsRead += shortsRead

		// WriteToWaveFile's numSamples is total shorts. Returns shorts written.
		// Only write the number of shorts actually read.
		okWritten := outputFile.WriteToWaveFile(buffer[:shortsRead], shortsRead)
		if okWritten == 0 {
			t.Fatalf("WriteToWaveFile returned false value: %d", okWritten)
		}
	}

	if totalShortsRead != expectedTotalShortsToRead {
		t.Errorf("Expected to read %d total shorts, but read %d. Dummy WAV format or C reader might be the issue.", expectedTotalShortsToRead, totalShortsRead)
	}

	t.Logf("Successfully read %d shorts.", totalShortsRead)
}
