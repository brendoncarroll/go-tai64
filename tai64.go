package tai64

import (
	"encoding/binary"
	"fmt"
	"time"
)

const nano = 1e9

type TAI64 uint64

func (t TAI64) String() string {
	return fmt.Sprintf("TAI64(%d)", t)
}

func (t TAI64) Marshal() (ret [8]byte) {
	binary.BigEndian.PutUint64(ret[:], uint64(t))
	return ret
}

func (t TAI64) TAI64N() TAI64N {
	return TAI64N{
		Seconds: uint64(t),
	}
}

func (t TAI64) TAI64NA() TAI64NA {
	return TAI64NA{
		Seconds: uint64(t),
	}
}

func Parse(x []byte) (TAI64, error) {
	if len(x) != 8 {
		return 0, errWrongLength("TAI64", len(x), 8)
	}
	return TAI64(binary.BigEndian.Uint64(x)), nil
}

type TAI64N struct {
	Seconds     uint64 `json:"s"`
	Nanoseconds uint32 `json:"ns"`
}

func (t TAI64N) String() string {
	return fmt.Sprintf("TAI64N(s=%d, ns=%d)", t.Seconds, t.Nanoseconds)
}

func (t TAI64N) Marshal() (ret [12]byte) {
	binary.BigEndian.PutUint64(ret[0:8], uint64(t.Seconds))
	binary.BigEndian.PutUint32(ret[8:12], uint32(t.Nanoseconds))
	return ret
}

// TAI64 returns a TAI64, truncating the nanoseconds.
func (t TAI64N) TAI64() TAI64 {
	return TAI64(t.Seconds)
}

func (t TAI64N) TAI64NA() TAI64NA {
	return TAI64NA{
		Seconds:     t.Seconds,
		Nanoseconds: t.Nanoseconds,
	}
}

func ParseN(data []byte) (TAI64N, error) {
	if len(data) != 16 {
		return TAI64N{}, errWrongLength("TAI64N", len(data), 12)
	}
	seconds := binary.BigEndian.Uint64(data[0:8])
	nanos := binary.BigEndian.Uint32(data[8:12])
	if nanos >= nano {
		return TAI64N{}, fmt.Errorf("nanoseconds exceed 1e9: %d", nanos)
	}
	return TAI64N{
		Seconds:     seconds,
		Nanoseconds: nanos,
	}, nil
}

type TAI64NA struct {
	Seconds     uint64 `json:"s"`
	Nanoseconds uint32 `json:"ns"`
	Attoseconds uint32 `json:"as"`
}

func (t TAI64NA) String() string {
	return fmt.Sprintf("TAI64NA(s=%d, ns=%d, as=%d)", t.Seconds, t.Nanoseconds, t.Attoseconds)
}

func (t TAI64NA) Marshal() (ret [16]byte) {
	binary.BigEndian.PutUint64(ret[0:8], uint64(t.Seconds))
	binary.BigEndian.PutUint32(ret[8:12], uint32(t.Nanoseconds))
	binary.BigEndian.PutUint32(ret[12:16], uint32(t.Attoseconds))
	return ret
}

// TAI64N returns a TAI64N, truncating the attoseconds.
func (t TAI64NA) TAI64N() TAI64N {
	return TAI64N{
		Seconds:     t.Seconds,
		Nanoseconds: t.Nanoseconds,
	}
}

func ParseNA(data []byte) (TAI64NA, error) {
	if len(data) != 16 {
		return TAI64NA{}, errWrongLength("TAI64NA", len(data), 16)
	}
	seconds := binary.BigEndian.Uint64(data[0:8])
	nanos := binary.BigEndian.Uint32(data[8:12])
	attos := binary.BigEndian.Uint32(data[12:16])
	if nanos >= nano {
		return TAI64NA{}, fmt.Errorf("nanoseconds exceed 1e9: %d", nanos)
	}
	if attos >= nano {
		return TAI64NA{}, fmt.Errorf("attoseconds exceed 1e9: %d", attos)
	}
	return TAI64NA{
		Seconds:     seconds,
		Nanoseconds: nanos,
		Attoseconds: attos,
	}, nil
}

// GoTime returns a time.Time representing UTC time.
func (t TAI64) GoTime() time.Time {
	seconds := UTCFromTAI(taiSeconds(t))
	return time.Unix(seconds, 0).UTC()
}

// GoTime returns a time.Time representing UTC time.
func (t TAI64N) GoTime() time.Time {
	seconds := UTCFromTAI(taiSeconds(TAI64(t.Seconds)))
	return time.Unix(seconds, int64(t.Nanoseconds)).UTC()
}

// FromGoTime returns a TAI64N from a time.Time
// The time will be converted to UTC and then into TAI, and then finally into a TAI64N
func FromGoTime(x time.Time) TAI64N {
	nanos := x.UTC().UnixNano()
	seconds := nanos / nano
	nanos = nanos % nano
	return TAI64N{
		Seconds:     uint64(tai64(TAIFromUTC(seconds))),
		Nanoseconds: uint32(nanos),
	}
}

// Now is equivalent to FromGoTime(time.Now())
func Now() TAI64N {
	return FromGoTime(time.Now())
}

const tai64Offset = 1 << 62

func tai64(taiSecs int64) TAI64 {
	return TAI64(taiSecs + tai64Offset)
}

func taiSeconds(x TAI64) int64 {
	return int64(x) - tai64Offset
}
