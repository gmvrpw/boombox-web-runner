package runner

import (
	"bytes"

	rodstream "github.com/navicstein/rod-stream"
	"github.com/remko/go-mkvparse"
)

var (
	Channels     = 2
	SampleRate   = 48000
	FrameSize    = 960
	PcmSize      = FrameSize * Channels
	MaxDataBytes = FrameSize * Channels * 2
)

func (r *RunnerService) pipe(raw chan string, opus chan []byte) {
	b := make(chan []byte, 100)
	blocks := make(chan OpusBlock, 100)
	pcm := make(chan []int16, 100)

	go r.decodeBase64(raw, b)
	go r.extractOpusBlocks(b, blocks)
	go r.opusBlockToPCM(blocks, pcm)
	go r.pcmToOpus(pcm, opus)
}

func (r *RunnerService) decodeBase64(b64 chan string, bs chan []byte) {
	defer close(bs)

	for str := range b64 {
		bs <- rodstream.Parseb64(str)
	}
}

func (r *RunnerService) extractOpusBlocks(bs chan []byte, blocks chan OpusBlock) {
	defer close(blocks)

	h := AudioStreamHandler{blocks: blocks}

	var batch []byte
	for b := range bs {
		if len(b) == 1 {
			batch = b
			continue
		}
		batch = append(batch, b[:len(b)-1]...)
		mkvparse.Parse(bytes.NewReader(batch), &h)
		batch = []byte{b[len(b)-1]}
	}
}

func (r *RunnerService) opusBlockToPCM(blocks chan OpusBlock, pcm chan []int16) {
	defer close(pcm)

	for b := range blocks {
		if b.Duration == 0 {
			continue
		}
		d, err := r.decoder.Decode(b.Data, SampleRate/1000*b.Duration, false)
		if err != nil {
			// TODO: log this
			continue
		}

		pcm <- d
	}
}

func (r *RunnerService) pcmToOpus(pcm chan []int16, opus chan []byte) {
	defer close(opus)

	// Rest of PCM
	rest := []int16{}
	// Length of complementary of rest
	comp := 0

	for b := range pcm {
		comp = PcmSize - len(rest)
		if len(b) < comp {
			rest = append(rest, b...)
		}

		// Start of PCM frame
		s := 0
		for e := comp; e < len(b); e += PcmSize {
			o, err := r.encoder.Encode(append(rest, b[s:e]...), FrameSize, MaxDataBytes)
			if err != nil {
				// TODO: log error
				return
			}
			if len(rest) != 0 {
				rest = []int16{}
			}

			opus <- o
			s = e
		}
		rest = b[s:]
	}
}
