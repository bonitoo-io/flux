// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbsemantic

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type DateTimeLiteral struct {
	_tab flatbuffers.Table
}

func GetRootAsDateTimeLiteral(buf []byte, offset flatbuffers.UOffsetT) *DateTimeLiteral {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DateTimeLiteral{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *DateTimeLiteral) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DateTimeLiteral) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DateTimeLiteral) Loc(obj *SourceLocation) *SourceLocation {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(SourceLocation)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *DateTimeLiteral) Value(obj *Time) *Time {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Time)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func DateTimeLiteralStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func DateTimeLiteralAddLoc(builder *flatbuffers.Builder, loc flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(loc), 0)
}
func DateTimeLiteralAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(value), 0)
}
func DateTimeLiteralEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}