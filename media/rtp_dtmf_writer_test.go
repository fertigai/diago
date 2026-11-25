// SPDX-License-Identifier: MPL-2.0

package media

import (
	"testing"

	"github.com/pion/rtp"
	"github.com/stretchr/testify/require"
)

type rtpCapture []rtp.Packet

func (c *rtpCapture) WriteRTP(p *rtp.Packet) error {
	*c = append(*c, *p)
	return nil
}

func TestRTPDTMFWriterBackToBack(t *testing.T) {
	var pkts rtpCapture
	pw := NewRTPPacketWriter(&pkts, CodecAudioUlaw)
	w := NewRTPDTMFWriter(CodecTelephoneEvent8000, pw, pw)
	require.NoError(t, w.WriteDTMF('1'))
	require.NoError(t, w.WriteDTMF('2'))
	require.Len(t, pkts, 14)
	first, second := pkts[0].Timestamp, pkts[7].Timestamp
	require.Equal(t, first, pkts[6].Timestamp, "one event keeps its start timestamp")
	require.Equal(t, first+7*160, second, "the next event starts where this one ends")
}
