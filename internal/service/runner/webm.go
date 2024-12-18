package runner

import (
	"encoding/binary"

	"github.com/remko/go-mkvparse"
)

type OpusBlock struct {
	Data     []byte
	Duration int
}

type AudioStreamHandler struct {
	mkvparse.DefaultHandler

	// Timecodes
	segment int64
	to      int64
	from    int64

	// Out channel
	blocks chan OpusBlock
}

func (h *AudioStreamHandler) HandleInteger(
	id mkvparse.ElementID,
	time int64,
	_ mkvparse.ElementInfo,
) error {
	if id == mkvparse.TimecodeElement {
		h.segment = time
	}
	return nil
}

func (h *AudioStreamHandler) HandleBinary(
	id mkvparse.ElementID,
	data []byte,
	info mkvparse.ElementInfo,
) error {
	switch id {
	case mkvparse.SimpleBlockElement:
		h.to = h.segment + int64(binary.BigEndian.Uint16(data[1:3]))
		duration := int(h.to - h.from)

		// NOTE: It just work idk why
		// TODO: Get understanding and fix it
		if duration < 60 {
			h.blocks <- OpusBlock{Duration: 60, Data: data[4:]}
		} else {
			h.blocks <- OpusBlock{Duration: duration, Data: data[4:]}
		}

		h.from = h.to
	}
	return nil
}
