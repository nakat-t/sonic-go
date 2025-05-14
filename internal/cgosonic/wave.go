package cgosonic

/*
#cgo CFLAGS: -Wall -Wno-unused-function -g -ansi -fPIC -pthread -I${SRCDIR}
#include <stdlib.h>
#include "wave.h"
*/
import "C"
import (
	"errors"
	"os"
	"unsafe"
)

// WaveFile represents a WAVE file
type WaveFile struct {
	file C.waveFile
}

// OpenInputWaveFile opens an input WAVE file
func OpenInputWaveFile(fileName string) (*WaveFile, int, int, error) {
	// openInputWaveFile outputs to stderr if file open fails.
	// So, check here to prevent output.
	f, err := os.Open(fileName)
	if err != nil {
		return nil, 0, 0, err
	}
	f.Close()

	var sampleRate C.int
	var numChannels C.int
	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	file := C.openInputWaveFile(cFileName, &sampleRate, &numChannels)
	if file == nil {
		return nil, 0, 0, errors.New("failed to open input wave file")
	}
	return &WaveFile{file: file}, int(sampleRate), int(numChannels), nil
}

// OpenOutputWaveFile opens an output WAVE file
func OpenOutputWaveFile(fileName string, sampleRate int, numChannels int) (*WaveFile, error) {
	// openOutputWaveFile outputs to stderr if file open fails.
	// So, check here to prevent output.
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	f.Close()

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	file := C.openOutputWaveFile(cFileName, C.int(sampleRate), C.int(numChannels))
	if file == nil {
		return nil, errors.New("failed to open output wave file")
	}
	return &WaveFile{file: file}, nil
}

// CloseWaveFile closes a WAVE file
func (w *WaveFile) CloseWaveFile() int {
	if w.file == nil {
		// closeWaveFile returns 1 if success, 0 if fail.
		return 1
	}
	result := C.closeWaveFile(w.file)
	w.file = nil
	return int(result)
}

// ReadFromWaveFile reads samples from a WAVE file
func (w *WaveFile) ReadFromWaveFile(buffer []int16, maxSamples int) int {
	return int(C.readFromWaveFile(w.file, (*C.short)(unsafe.Pointer(&buffer[0])), C.int(maxSamples)))
}

// WriteToWaveFile writes samples to a WAVE file
func (w *WaveFile) WriteToWaveFile(buffer []int16, numSamples int) int {
	return int(C.writeToWaveFile(w.file, (*C.short)(unsafe.Pointer(&buffer[0])), C.int(numSamples)))
}
