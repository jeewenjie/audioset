package audioset

import (
	"bytes"
	"io"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/unixpickle/wav"
)

func TestReadAndMix(t *testing.T) {
	for _, numChan := range []int{1, 2, 3} {
		t.Run("Chan"+strconv.Itoa(numChan), func(t *testing.T) {
			sound := wav.NewPCM16Sound(numChan, 22050)
			initSamples := make([]wav.Sample, 22050*numChan*10)
			for i := range initSamples {
				initSamples[i] = wav.Sample(rand.Float64()*2 - 1)
			}
			sound.SetSamples(initSamples)

			var buf bytes.Buffer
			if err := sound.Write(&buf); err != nil {
				t.Fatal(err)
			}

			reader := bytes.NewReader(buf.Bytes())
			data, err := readAndMix(reader)
			if err != nil {
				t.Fatal(err)
			}
			for i, actual := range data {
				var expected float64
				for j := 0; j < numChan; j++ {
					expected += float64(initSamples[i*numChan+j]) / float64(numChan)
				}
				if math.Abs(actual-expected) > 1e-2 {
					t.Errorf("bad data: got %v expected %v", actual, expected)
					break
				}
			}
		})
	}
}

func BenchmarkReadAndMix(b *testing.B) {
	sound := wav.NewPCM16Sound(2, 22050)
	sound.SetSamples(make([]wav.Sample, 22050*2*10))
	var buf bytes.Buffer
	if err := sound.Write(&buf); err != nil {
		b.Fatal(err)
	}

	reader := bytes.NewReader(buf.Bytes())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		readAndMix(reader)
	}
}
