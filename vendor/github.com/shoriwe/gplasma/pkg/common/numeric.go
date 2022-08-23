package common

import (
	"encoding/binary"
	"math"
)

type (
	integer interface {
		int | int8 | int16 | int32 | int64
	}
	float interface {
		float32 | float64
	}
)

func IntToBytes[T integer](i T) []byte {
	var bytes [8]byte
	binary.BigEndian.PutUint64(bytes[:], uint64(i))
	return bytes[:]
}

func BytesToInt(i []byte) int64 {
	switch len(i) {
	case 1:
		return int64(i[0])
	case 2:
		return int64(binary.BigEndian.Uint16(i))
	case 4:
		return int64(binary.BigEndian.Uint32(i))
	case 8:
		return int64(binary.BigEndian.Uint64(i))
	default:
		panic("invalid integer length")
	}
}

func FloatToBytes[T float](f T) []byte {
	var bytes [8]byte
	u := math.Float64bits(float64(f))
	binary.BigEndian.PutUint64(bytes[:], u)
	return bytes[:]
}

func BytesToFloat(i []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(i))
}
