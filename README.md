# sonic-go

sonic-go is go bindings for [Sonic](https://android.googlesource.com/platform/external/sonic/+/master/doc/index.md). Sonic is a algorithm for changing speech speed, uniquely optimized for over 2X acceleration.

## Overview

* Interface compatible with Go's standard `io.Writer`
* Sonic allows you to change the speed of the audio. It is optimized for speeds of 2x or more.
* Pitch and volume can be changed at the same time.
* Supported wav audio format: LPCM(16bit signed) and IEEE float(32bit float)
* Support multi channels: 1(mono) to 32ch

## Installation

```bash
go get github.com/nakat-t/sonic-go
```

## Usage

The core of the sonic package is the `sonic.Transformer` class. This object wraps an `io.Writer` object and creates another `io.Writer` object that provides functions to modify the volume, pitch, speed, etc. of audio data.

Here is an example of 2.5x speed up and 1.5x volume increase for linear PCM audio:

```go
import "github.com/nakat-t/sonic-go"

...

// Input audio data
const sampleRate = 44100 // 44.1kHz
var audioData []byte = readAudioDataSomeWay() // 16bit signed, linear PCM 

// Output audio buffer
var outAudioData bytes.Buffer

// Create a Sonic transformer
trf, err := sonic.NewTransformer(outAudioData, sampleRate, sonic.AudioFormatPCM,
	sonic.WithSpeed(2.5),
	sonic.WithVolume(1.5),
)

// Write audio data to Sonic transformer. The audio with the volume 
// and speed changed is written to outAudioData via the transformer.
io.Copy(trf, bytes.NewBuffer(audioData))

// NOTE: After writing all input data, be sure to execute `Flush()` or
// `Close()` to write all audio data remaining in the internal buffer.
trf.Close()

// The converted audio data is output to `outAudioData`.
playAudioDataSomeWay(outAudioData)
```

## License

sonic-go is provided under the [Apache-2.0 license](./LICENSE) (same as sonic).
