#!/bin/bash
set -euo pipefail

script_dir="$(dirname "$(realpath "$0")")"
sonic_dir="$script_dir/../submodules/sonic"
sonic_bin="$sonic_dir/sonic"
testdata_dir="$script_dir/../test/testdata/reference"

original_wav="$testdata_dir/original/common_voice_en_1dcef00e46910f33.wav"

# make sonic binary for reference testdata generation
make -C "$sonic_dir" clean sonic

volume_test_values=(
    0.01
    0.5
    1.0
    2.0
    100.0
)

# Generate test data for volume tests
for volume in "${volume_test_values[@]}"; do
    testdata_file="$testdata_dir/volume_$volume.wav"
    "$sonic_bin" -v "$volume" "$original_wav" "$testdata_file"
done

speed_test_values=(
    0.05
    0.5
    1.0
    2.0
    20.0
)

# Generate test data for speed tests
for speed in "${speed_test_values[@]}"; do
    testdata_file="$testdata_dir/speed_$speed.wav"
    "$sonic_bin" -s "$speed" "$original_wav" "$testdata_file"
done

pitch_test_values=(
    0.05
    0.5
    1.0
    2.0
    20.0
)

# Generate test data for pitch tests
for pitch in "${pitch_test_values[@]}"; do
    testdata_file="$testdata_dir/pitch_$pitch.wav"
    "$sonic_bin" -p "$pitch" "$original_wav" "$testdata_file"
done

# Generate test data for quality tests
testdata_file="$testdata_dir/quality_on.wav"
"$sonic_bin" -q "$original_wav" "$testdata_file"

# clean sonic binary
make -C "$sonic_dir" clean
