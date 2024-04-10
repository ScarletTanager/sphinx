package probability_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProbability(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Probability Suite")
}
