package client_test

import (
	"fmt"

	. "github.com/Manan007224/sigo_client/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/ztrue/tracerr"
)

var (
	ex  *Executor
	err error
)

var _ = Describe("Executor", func() {
	Specify("Positive Tests", func() {
		ex = NewExecutor()
		testFunc := func(t string) error {
			return tracerr.Wrap(errors.New("testFunc"))
		}
		ex.Register(testFunc, "testFunc", "test")
		err, jErr := ex.Perform("testFunc", "test", 2, 2)
		Expect(err).ShouldNot(HaveOccurred())
		fmt.Println(jErr)
	})

	Specify("Invalid args provided", func() {
		ex = NewExecutor()
		testFunc := func(t string) error {
			return tracerr.Wrap(errors.New("testFunc"))
		}
		ex.Register(testFunc, "testFunc", "test")
		err, _ := ex.Perform("testFunc", 1, 2, 2)
		Expect(err).Should(HaveOccurred())
	})

	Specify("Invalid funcName", func() {
		ex = NewExecutor()
		testFunc := func(t string) error {
			return tracerr.Wrap(errors.New("testFunc"))
		}
		ex.Register(testFunc, "testFunc", "test")
		err, _ := ex.Perform("testFun", "test", 2, 2)
		Expect(err).Should(HaveOccurred())
	})
})
