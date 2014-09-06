//
// Copyright 2014 Hong Miao. All Rights Reserved.
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


package amf

import (
	"errors"
	"log"
	"github.com/oikomi/gortmpserver/server/util"
)

type Amf0Codec struct {
	buf *AmfBuffer
}

func NewAmf0Codec(b []byte) (*Amf0Codec) {
	return &Amf0Codec{
		buf : NewAmfBuffer(b),
	}
}

func (self *Amf0Codec)ReadString() (string, error) {
	var err error
	if marker := self.buf.b[0]; marker != AMF0_String {
		err = errors.New("amf0 string marker invalid")
		log.Fatalln(err.Error())
		return "", err
	}

	len := util.ReadUint16(self.buf.b[1:3])
	
	////log.Println(len)	
	v := string(self.buf.b[3:len+3])
	
	return v, nil
}

type Amf0CmdObject struct {
	marker byte
	properties map[string]interface{}
}
func NewAmf0CmdObject() (*Amf0CmdObject) {
	return &Amf0CmdObject {
		marker : AMF0_Object,
		properties : make(map[string]interface{}),
	}
}


