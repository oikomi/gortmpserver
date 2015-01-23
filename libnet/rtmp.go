package libnet


type rtmpProtocol struct {
	n             int
	bo            binary.ByteOrder
}

func newrtmpProtocol(n int, byteOrder binary.ByteOrder) *rtmpProtocol {
	return &rtmpProtocol {
		
	}
}


func (p *rtmpProtocol) PrepareOutBuffer(buffer *OutBuffer, size int) {
	buffer.Prepare(size)
	buffer.Data = buffer.Data[:p.n]
}

func (p *rtmpProtocol) Write(writer io.Writer, packet *OutBuffer) error {
	if p.MaxPacketSize > 0 && len(packet.Data) > p.MaxPacketSize {
		return PacketTooLargeError
	}
	p.encodeHead(packet.Data)
	if _, err := writer.Write(packet.Data); err != nil {
		return err
	}
	return nil
}

func (p *rtmpProtocol) Read(reader io.Reader, buffer *InBuffer) error {
	// head
	buffer.Prepare(p.n)
	if _, err := io.ReadFull(reader, buffer.Data); err != nil {
		return err
	}
	size := p.decodeHead(buffer.Data)
	if p.MaxPacketSize > 0 && size > p.MaxPacketSize {
		return PacketTooLargeError
	}
	// body
	buffer.Prepare(size)
	if size == 0 {
		return nil
	}
	if _, err := io.ReadFull(reader, buffer.Data); err != nil {
		return err
	}
	return nil
}