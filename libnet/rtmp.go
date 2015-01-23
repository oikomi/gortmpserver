package libnet

import (
	"encoding/binary"
	"io"
)

const (
	RTMP_HANDSHAKE_SERVER_RECV_C0C1      =  1
	RTMP_HANDSHAKE_SERVER_SEND_S0S1S2    =  2
	RTMP_HANDSHAKE_SERVER_RECV_C2        =  3
)

const (
	RTMP_HANDSHAKE_C0_LENGTH  =  1
	RTMP_HANDSHAKE_C1_LENGTH  =  1536
	RTMP_HANDSHAKE_C2_LENGTH  =  1536
	RTMP_HANDSHAKE_S0_LENGTH  =  1
	RTMP_HANDSHAKE_S1_LENGTH  =  1536
	RTMP_HANDSHAKE_S2_LENGTH  =  1536
)

type rtmpProtocol struct {
	stage         uint32
	bo            binary.ByteOrder
}

func NewrtmpProtocol(stage uint32, byteOrder binary.ByteOrder) *rtmpProtocol {
	return &rtmpProtocol {
		stage : stage,
		bo    : byteOrder,
	}
}

func (p *rtmpProtocol) New(v interface{}) ProtocolState {
	return p
}

func (p *rtmpProtocol) PrepareOutBuffer(buffer *OutBuffer, size int) {
	buffer.Prepare(size)
	//buffer.Data = buffer.Data[:p.n]
}

func (p *rtmpProtocol) Write(writer io.Writer, packet *OutBuffer) error {
	if p.stage == RTMP_HANDSHAKE_SERVER_SEND_S0S1S2 {
		if _, err := writer.Write(packet.Data); err != nil {
			return err
		}
		p.stage++
	}

	return nil
}

func (p *rtmpProtocol) Read(reader io.Reader, buffer *InBuffer) error {
	switch p.stage {
	case RTMP_HANDSHAKE_SERVER_RECV_C0C1:
		buffer.Prepare(RTMP_HANDSHAKE_C0_LENGTH + RTMP_HANDSHAKE_C1_LENGTH)
		if _, err := io.ReadFull(reader, buffer.Data); err != nil {
			return err
		}
		p.stage++
	case RTMP_HANDSHAKE_SERVER_RECV_C2:
		buffer.Prepare(RTMP_HANDSHAKE_C2_LENGTH)
		if _, err := io.ReadFull(reader, buffer.Data); err != nil {
			return err
		}
		p.stage++
	}


	return nil
}