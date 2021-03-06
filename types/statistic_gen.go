package types

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ConfirmTime) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if zb0001 != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zb0001}
		return
	}
	z.SeqHeight, err = dc.ReadUint64()
	if err != nil {
		return
	}
	z.TxNum, err = dc.ReadUint64()
	if err != nil {
		return
	}
	z.ConfirmTime, err = dc.ReadString()
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z ConfirmTime) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 3
	err = en.Append(0x93)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.SeqHeight)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.TxNum)
	if err != nil {
		return
	}
	err = en.WriteString(z.ConfirmTime)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z ConfirmTime) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 3
	o = append(o, 0x93)
	o = msgp.AppendUint64(o, z.SeqHeight)
	o = msgp.AppendUint64(o, z.TxNum)
	o = msgp.AppendString(o, z.ConfirmTime)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ConfirmTime) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if zb0001 != 3 {
		err = msgp.ArrayError{Wanted: 3, Got: zb0001}
		return
	}
	z.SeqHeight, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		return
	}
	z.TxNum, bts, err = msgp.ReadUint64Bytes(bts)
	if err != nil {
		return
	}
	z.ConfirmTime, bts, err = msgp.ReadStringBytes(bts)
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z ConfirmTime) Msgsize() (s int) {
	s = 1 + msgp.Uint64Size + msgp.Uint64Size + msgp.StringPrefixSize + len(z.ConfirmTime)
	return
}
