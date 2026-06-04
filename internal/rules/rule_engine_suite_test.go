package rules_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRuleEngine(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RuleEngine Suite")
}
