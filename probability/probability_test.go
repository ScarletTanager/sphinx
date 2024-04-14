package probability_test

import (
	"math/rand"
	"slices"

	"github.com/ScarletTanager/sphinx/probability"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Probability", func() {
	Describe("Discretize", func() {
		var (
			vals     []float64
			valCount int
			cfg      probability.DiscretizationConfig
		)

		BeforeEach(func() {
			cfg = probability.DiscretizationConfig{}
			valCount = 100
		})

		JustBeforeEach(func() {
			vals = make([]float64, valCount)
			for i := 0; i < valCount; i++ {
				vals[i] = rand.Float64()
			}
		})

		When("Using only the defaults", func() {
			It("Returns ten intervals", func() {
				intervals := probability.Discretize(vals, cfg)
				Expect(intervals).To(HaveLen(10))
			})

			It("Returns intervals of equal size", func() {
				intervals := probability.Discretize(vals, cfg)
				Expect(intervalsHaveEqualSize(intervals)).To(BeTrue())
			})

			It("Returns nonoverlapping intervals in the correct order", func() {
				intervals := probability.Discretize(vals, cfg)
				for idx := 1; idx < len(intervals); idx++ {
					Expect(intervals[idx].Lower).To(BeNumerically(">=", (intervals[idx-1].Lower + intervals[idx-1].Size)))
				}
			})

			It("Returns intervals spanning the entire range of values", func() {
				intervals := probability.Discretize(vals, cfg)
				slices.Sort(vals)
				Expect(intervals[0].Lower).To(Equal(vals[0]))
				Expect(intervals[len(intervals)-1].Upper).To(Equal(vals[len(vals)-1]))
			})
		})

		When("Only interval count is specified", func() {
			BeforeEach(func() {
				cfg.Intervals = 5
			})

			It("Returns the specified number of intervals", func() {
				intervals := probability.Discretize(vals, cfg)
				Expect(intervals).To(HaveLen(cfg.Intervals))
			})

			It("Uses IntervalEqualSize", func() {
				intervals := probability.Discretize(vals, cfg)
				Expect(intervalsHaveEqualSize(intervals)).To(BeTrue())
			})
		})

		When("Using equal distribution", func() {
			BeforeEach(func() {
				cfg.Method = probability.IntervalEqualDistribution
			})

			When("No interval count is specified", func() {
				It("Returns ten intervals", func() {
					intervals := probability.Discretize(vals, cfg)
					Expect(intervals).To(HaveLen(10))
				})
			})

			It("Returns intervals with an equal number of values per interval", func() {
				intervals := probability.Discretize(vals, cfg)
				expectedLen := intervalLen(intervals[0], vals)
				for _, i := range intervals {
					Expect(intervalLen(i, vals)).To(Equal(expectedLen))
				}
			})

			When("An interval count is specified", func() {
				BeforeEach(func() {
					cfg.Intervals = 5
				})

				It("Returns the specified number of intervals", func() {
					intervals := probability.Discretize(vals, cfg)
					Expect(intervals).To(HaveLen(cfg.Intervals))
				})

				It("Returns intervals with an equal number of values per interval", func() {
					intervals := probability.Discretize(vals, cfg)
					expectedLen := intervalLen(intervals[0], vals)
					for _, i := range intervals {
						Expect(intervalLen(i, vals)).To(Equal(expectedLen))
					}
				})

				When("The number of values is not evenly divisible by the specified interval count", func() {
					BeforeEach(func() {
						cfg.Intervals = 7
						Expect(valCount % cfg.Intervals).NotTo(Equal(0))
					})

					It("Returns the specified number of intervals", func() {
						intervals := probability.Discretize(vals, cfg)
						Expect(intervals).To(HaveLen(cfg.Intervals))
					})

					It("Assigns an equal number of values to all but the last interval", func() {
						intervals := probability.Discretize(vals, cfg)
						expectedLen := intervalLen(intervals[0], vals)
						for _, i := range intervals[:cfg.Intervals-1] {
							Expect(intervalLen(i, vals)).To(Equal(expectedLen))
						}
					})

					It("Assigns the remainder ((valCount / intervalCount) + (valCount % intervalCount)) to the last interval", func() {
						intervals := probability.Discretize(vals, cfg)
						expectedLen := intervalLen(intervals[0], vals) + (valCount % cfg.Intervals)
						Expect(intervalLen(intervals[cfg.Intervals-1], vals)).To(Equal(expectedLen))
					})
				})
			})
		})

		When("Using equal size", func() {
			BeforeEach(func() {
				cfg.Method = probability.IntervalEqualSize
			})

			When("No interval count is specified", func() {
				It("Returns ten intervals", func() {
					intervals := probability.Discretize(vals, cfg)
					Expect(intervals).To(HaveLen(10))
				})

				It("Returns intervals of equal size", func() {
					intervals := probability.Discretize(vals, cfg)
					Expect(intervalsHaveEqualSize(intervals)).To(BeTrue())
				})
			})

			When("An interval count is specified", func() {
				BeforeEach(func() {
					cfg.Intervals = 5
				})

				It("Returns the specified number of intervals", func() {
					intervals := probability.Discretize(vals, cfg)
					Expect(intervals).To(HaveLen(cfg.Intervals))
				})

				It("Returns intervals of equal size", func() {
					intervals := probability.Discretize(vals, cfg)
					Expect(intervalsHaveEqualSize(intervals)).To(BeTrue())
				})
			})
		})
	})

	Describe("Interval", func() {
		var (
			i                            *probability.Interval
			l, u, s                      float64
			includesLower, includesUpper bool
		)

		BeforeEach(func() {
			l = 10.0
			u = 0.0
			s = 0.0
			includesLower = false
			includesUpper = false
		})

		JustBeforeEach(func() {
			i = &probability.Interval{
				Lower:         l,
				Upper:         u,
				Size:          s,
				IncludesLower: includesLower,
				IncludesUpper: includesUpper,
			}
		})

		Describe("Contains", func() {
			When("The interval has size specified (discretization using equal size)", func() {
				BeforeEach(func() {
					s = 5.0
					includesLower = true
				})

				When("The val is contained within the interval", func() {
					It("Returns true", func() {
						Expect(i.Contains(i.Lower + (i.Size / 2.0))).To(BeTrue())
					})
				})

				When("The val is not contained within the interval", func() {
					It("Returns false", func() {
						Expect(i.Contains(i.Lower + (i.Size * 2.0))).To(BeFalse())
					})
				})

				When("The val equals the lower limit", func() {
					It("Returns true", func() {
						Expect(i.Contains(i.Lower)).To(BeTrue())
					})
				})

				When("The val equals the upper limit (lower + size)", func() {
					It("Returns false", func() {
						Expect(i.Contains(i.Lower + i.Size)).To(BeFalse())
					})

					When("This is the last interval in the range", func() {
						BeforeEach(func() {
							includesUpper = true
						})

						It("Returns true", func() {
							Expect(i.Contains(i.Lower + i.Size)).To(BeTrue())
						})
					})
				})
			})

			When("The interval has upper specified (discretization using equal distribution)", func() {
				BeforeEach(func() {
					u = 20.0
					includesUpper = true
				})

				JustBeforeEach(func() {
					Expect(i.Size).To(Equal(0.0))
				})

				When("The val is contained within the interval", func() {
					It("Returns true", func() {
						Expect(i.Contains(i.Lower + ((i.Upper - i.Lower) / 2.0))).To(BeTrue())
					})
				})

				When("The val is not contained within the interval", func() {
					It("Returns false", func() {
						Expect(i.Contains(i.Upper + i.Lower)).To(BeFalse())
					})
				})

				When("The val equals the lower bound", func() {
					It("Returns true", func() {
						Expect(i.Contains(i.Lower)).To(BeTrue())
					})
				})

				When("The Val equals the upper bound", func() {
					It("Returns true", func() {
						Expect(i.Contains(i.Upper)).To(BeTrue())
					})
				})
			})
		})
	})
})

func intervalsHaveEqualSize(intervals []probability.Interval) bool {
	size := intervals[0].Size
	for _, interval := range intervals[:len(intervals)-1] {
		if interval.Size != size {
			return false
		}
	}

	return true
}

func intervalsHaveEqualLen(intervals []probability.Interval, vals []float64) bool {
	slices.Sort(vals)

	// Take the first interval, determine the number of values enclosed by the closed interval
	// [Lower, Upper]
	first := intervals[0].Lower
	Expect(first).To(Equal(vals[0]))
	last := intervals[0].Upper

	expectedLen := 0
	for _, v := range vals {
		if v <= last {
			expectedLen++
		} else {
			break
		}
	}

	for _, i := range intervals[1:] {
		if intervalLen(i, vals) != expectedLen {
			return false
		}
	}

	return true
}

func intervalLen(i probability.Interval, vals []float64) int {
	var thisLen int

	for _, v := range vals {
		if v >= i.Lower {
			if v <= i.Upper {
				thisLen++
			} else {
				break
			}
		}
	}

	return thisLen
}
