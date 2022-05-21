package tai64

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnixEpoch(t *testing.T) {
	x := UnixEpoch.GoTime()
	t.Log(x.String())

	require.Equal(t, 1970, x.Year())
	require.Equal(t, time.January, x.Month())
	require.Equal(t, 1, x.Day())
	require.Equal(t, 0, x.Hour())
	require.Equal(t, 0, x.Minute())
	require.Equal(t, 0, x.Second())
}

func TestMarshalParse(t *testing.T) {
	rng := rand.New(rand.NewSource(0))
	for i := 0; i < 100; i++ {
		x := TAI64((rng.Int63n(1 << 62)))
		data := x.Marshal()
		y, err := Parse(data[:])
		require.NoError(t, err)
		require.Equal(t, x, y)
	}
	for i := 0; i < 100; i++ {
		x := TAI64N{
			Seconds:     uint64(rng.Int63n(1 << 62)),
			Nanoseconds: uint32(rng.Int31n(nano)),
		}
		data := x.Marshal()
		y, err := ParseN(data[:])
		require.NoError(t, err)
		require.Equal(t, x, y)
	}
	for i := 0; i < 100; i++ {
		x := TAI64NA{
			Seconds:     uint64(rng.Int63n(1 << 62)),
			Nanoseconds: uint32(rng.Int31n(nano)),
			Attoseconds: uint32(rng.Int31n(nano)),
		}
		data := x.Marshal()
		y, err := ParseNA(data[:])
		require.NoError(t, err)
		require.Equal(t, x, y)
	}
}

func TestConvertUTC(t *testing.T) {
	type testCase struct {
		TAI int64
		UTC int64
	}
	// TAI, UTC
	tcs := []testCase{
		{10, 0},
		{11, 1},
		{-20, -30},
		{1637945456, 1637945419}, // Nov 2021, TAI = UTC + 37
	}
	for i, tc := range tcs {
		actualUTC := UTCFromTAI(tc.TAI)
		require.Equal(t, tc.UTC, actualUTC, "CASE %02d: UTCFromTAI is incorrect HAVE: %d, WANT: %d", i, actualUTC, tc.UTC)
		actualTAI := TAIFromUTC(tc.UTC)
		require.Equal(t, tc.TAI, actualTAI, "CASE %02d: TAIFromUTC is incorrect HAVE: %d, WANT: %d", i, actualTAI, tc.TAI)
	}
}
