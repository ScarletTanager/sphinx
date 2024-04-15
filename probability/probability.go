package probability

import (
	"slices"
)

type Interval struct {
	Lower         float64
	Upper         float64
	Size          float64
	IncludesLower bool
	IncludesUpper bool
}

// Contains returns true of the specified value is contained by the interval,
// false otherwise
func (i *Interval) Contains(val float64) bool {
	if i.Upper != 0.0 {
		// Either equal distribution or the last interval from equal size
		return val >= i.Lower && ((i.IncludesUpper && val <= i.Upper) || val < i.Upper)
	}

	// Interval was created using equal size
	if i.Size != 0.0 {
		if val >= i.Lower {
			return val < (i.Lower+i.Size) || (val == (i.Lower+i.Size) && i.IncludesUpper)
		}
	}

	// Either a screwed-up interval or the value is not contained
	return false
}

type Intervals []Interval

// IntervalForValue returns the index 0..(len(is) - 1) of the
// Interval containing the passed value
func (is Intervals) IntervalForValue(val float64) int {
	for idx, interval := range is {
		if interval.Contains(val) {
			return idx
		}
	}

	// I know, this isn't very idiomatic.  Should be an error, will fix later.
	return -1
}

// DiscretizationConfig controls the behavior of discretization of a continuous
// range of values.  Intervals is the number of intervals, Method determines
// how the range is subdivided, and IncludeUpperBound toggles whether each interval
// includes its upper bound (default of false means that only the last interval
// includes its upper bound, all others exclude it).
type DiscretizationConfig struct {
	Intervals         int
	Method            DiscretizationMethod
	IncludeUpperBound bool // Unused for now
}

type DiscretizationMethod int

const (
	IntervalEqualSize         DiscretizationMethod = iota // Every interval is the same size
	IntervalEqualDistribution DiscretizationMethod = iota // Every interval contains the same number of known values
	DefaultIntervalCount                           = 10
)

// Discretize converts a continuous (real-valued) range into a set of discrete intervals.
// vals is the set of known values within the range.  Pass an empty config object to
// get the default behavior of 10 intervals, each interval is of equal size, every interval
// but the last includes its lower bound and excludes its upper bound.  The last interval
// includes both bounds by default.
func Discretize(vals []float64, cfg DiscretizationConfig) []Interval {
	var (
		intervalCount int
	)

	if cfg.Intervals <= 0 {
		intervalCount = DefaultIntervalCount
	} else {
		intervalCount = cfg.Intervals
	}

	intervals := make([]Interval, intervalCount)

	slices.Sort(vals)

	// This assumes we know all the values we will ever see - this could be viewed
	// as a bug, but for now it's just a known limitation.  Why does this matter?
	// If you train a model with a a dataset, and you then try to classify based on
	// attribute values outside the lower/upper bounds of the original dataset...
	switch cfg.Method {
	case IntervalEqualSize:
		rangeSize := vals[len(vals)-1] - vals[0]
		intervalSize := rangeSize / float64(intervalCount)
		intervals[0] = Interval{
			Lower:         vals[0],
			Size:          intervalSize,
			IncludesLower: true,
			IncludesUpper: false,
		}

		for i := 1; i < intervalCount; i++ {
			intervals[i] = Interval{
				Lower:         intervals[i-1].Lower + intervalSize,
				Size:          intervalSize,
				IncludesLower: true,
				IncludesUpper: false,
			}
		}

		intervals[intervalCount-1].Upper = vals[len(vals)-1]
		intervals[intervalCount-1].IncludesUpper = true

	case IntervalEqualDistribution:
		// Using "length" to mean "count of values in interval" vs.
		// "size" for the size of the interval (b - a)
		intervalLen := len(vals) / intervalCount
		for i := 0; i < (intervalCount - 1); i++ {
			intervals[i] = Interval{
				Lower:         vals[i*intervalLen],
				Upper:         vals[((i+1)*intervalLen)-1],
				IncludesLower: true,
				IncludesUpper: true,
			}
		}

		intervals[intervalCount-1] = Interval{
			Lower:         vals[(intervalCount-1)*intervalLen],
			Upper:         vals[len(vals)-1],
			IncludesLower: true,
			IncludesUpper: true,
		}
	}

	return intervals
}
