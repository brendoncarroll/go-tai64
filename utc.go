package tai64

import "time"

// UnixEpoch as TAI64
const UnixEpoch = TAI64(tai64Offset + 10)

// TAIFromUTC applies leap seconds to seconds
// and returns the corresponding TAI time.
// unixSeconds is the number of seconds since the UNIX epoch
// The result IS NOT a TAI64
func TAIFromUTC(unixSeconds int64) int64 {
	return addLeapSeconds(leapSeconds1970, unixSeconds)
}

// UTCFromTAI removes leap seconds and
// returns the corresponding UTC time
// unixSeconds is the number of seconds since the UNIX epoch
func UTCFromTAI(unixSeconds int64) int64 {
	return removeLeapSeconds(leapSeconds1970, unixSeconds)
}

type leapSecond struct {
	UTC    int64
	Offset int64
}

var leapSeconds1970 = shiftLeapSeconds(leapSeconds1900)

func shiftLeapSeconds(x []leapSecond) (ret []leapSecond) {
	t1 := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	shift := int64(t2.Sub(t1).Seconds())
	ret = make([]leapSecond, len(x))
	for i := range x {
		ret[i] = x[i]
		ret[i].UTC -= shift
	}
	return ret
}

func addLeapSeconds(table []leapSecond, x int64) int64 {
	// TODO: binary search
	for i := range table {
		if i == len(table)-1 || table[i+1].UTC > x {
			return x + table[i].Offset
		}
	}
	panic("impossible")
}

func removeLeapSeconds(table []leapSecond, x int64) int64 {
	// TODO: binary search
	for i := range table {
		if i == len(table)-1 || table[i+1].UTC > x-table[i].Offset {
			return x - table[i].Offset
		}
	}
	panic("impossible")
}
