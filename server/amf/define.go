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


// AMF0 marker
const AMF0_Number = 0x00
const AMF0_Boolean = 0x01
const AMF0_String = 0x02
const AMF0_Object = 0x03
const AMF0_MovieClip = 0x04 // reserved, not supported
const AMF0_Null = 0x05
const AMF0_Undefined = 0x06
const AMF0_Reference = 0x07
const AMF0_EcmaArray = 0x08
const AMF0_ObjectEnd = 0x09
const AMF0_StrictArray = 0x0A
const AMF0_Date = 0x0B
const AMF0_LongString = 0x0C
const AMF0_UnSupported = 0x0D
const AMF0_RecordSet = 0x0E // reserved, not supported
const AMF0_XmlDocument = 0x0F
const AMF0_TypedObject = 0x10
// AVM+ object is the AMF3 object.
const AMF0_AVMplusObject = 0x11
// origin array whos data takes the same form as LengthValueBytes
const AMF0_OriginStrictArray = 0x20

// User defined
const AMF0_Invalid = 0x3F