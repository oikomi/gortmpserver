//
// Copyright 2014-2099 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import(
	"net"
	"bufio"
	"github.com/golang/glog"
)

type Session struct {
	conn             net.Conn
	br               *bufio.Reader
	bw               *bufio.Writer
	closed           bool
	inChunkStreams   map[uint32]*ChunkStream
	outChunkStreams  map[uint32]*ChunkStream
}

func NewSession(conn net.Conn) *Session {
	return &Session {
		conn             : conn,
		br               : bufio.NewReader(conn),
		bw               : bufio.NewWriter(conn),
		outChunkStreams  :  make(map[uint32]*ChunkStream),
		inChunkStreams   :  make(map[uint32]*ChunkStream),
	}
}


func (self *Session)bufRead(b []byte) error {
	_, err := self.br.Read(b)
	
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	return nil
}

func (self *Session)bufWrite(b []byte) error {
	_, err := self.bw.Write(b)
	
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	err = self.bw.Flush()
	if err != nil {
		glog.Error(err.Error())
		return err
	}
	
	return nil
}

func (self *Session)readLoop() error {
	for !self.closed {
		cs := NewChunkStream()
		cs.ReadChunkStream(self.br)
		self.inChunkStreams[cs.bh.chunkStreamID] = cs
	}
	return nil
}

func (self *Session)sendLoop() error {
	
	return nil
}
