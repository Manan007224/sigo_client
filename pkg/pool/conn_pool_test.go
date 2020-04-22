package pool_test

import (
	. "github.com/Manan007224/sigo_client/pkg/pool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

func closeConn(conn *grpc.ClientConn) {}

func createConn() (*grpc.ClientConn, error) {
	return &grpc.ClientConn{}, nil
}

func stubConn(cp *ConnPool) *ConnPool {
	cp.CreateHandler = createConn
	cp.CloseHandler = closeConn
	return cp
}

var _ = Describe("Pool", func() {
	Specify("Regular Pool", func() {
		cp := stubConn(NewConnPool(3, 6))
		seenConns := make(map[*grpc.ClientConn]bool)

		for i := 0; i < 5; i++ {
			conn, _ := cp.Get()
			seenConns[conn] = true
		}
		Expect(len(cp.FreeConns)).Should(Equal(0))

		for c := range seenConns {
			cp.Return(c)
		}

		Expect(len(cp.FreeConns)).Should(Equal(3))

		used := 0
		for i := 0; i < 5; i++ {
			conn, err := cp.Get()
			Expect(err).ShouldNot(HaveOccurred())
			if seenConns[conn] {
				used++
			}
		}

		Expect(used).Should(Equal(3))

		err := cp.Close()
		Expect(len(cp.FreeConns)).Should(Equal(0))
		Expect(err).ShouldNot(HaveOccurred())

		_, err = cp.Get()
		Expect(err).Should(HaveOccurred())
	})

	Specify("ConnPool Full", func() {
		cp := stubConn(NewConnPool(3, 6))
		conns := make([]*grpc.ClientConn, 6)
		for i := 0; i < 6; i++ {
			c, _ := cp.Get()
			conns[i] = c
		}
		_, err := cp.Get()
		Expect(err).Should(HaveOccurred())

		cp.Return(conns[0])
		_, err = cp.Get()
		Expect(err).ShouldNot(HaveOccurred())
	})
})
