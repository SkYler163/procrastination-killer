package signaller

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/pkg/errors"
)

// Signaller signaller struct.
type Signaller struct {
	buffer *beep.Buffer
}

// NewSignaller creates an instance of signaller.
func NewSignaller(filepath string) (*Signaller, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode mp3")
	}

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return nil, errors.Wrap(err, "failed init speaker")
	}

	volume := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   -0.5,
		Silent:   false,
	}

	buffer := beep.NewBuffer(format)
	buffer.Append(volume)

	err = streamer.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close streamer")
	}

	return &Signaller{buffer: buffer}, nil
}

// Signal makes a sound to mark the end of the period.
func (s *Signaller) Signal() {
	speaker.Play(s.buffer.Streamer(0, s.buffer.Len()))
}
