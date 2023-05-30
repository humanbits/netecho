package sender

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestTCPSender(t *testing.T) {

	host, port, getMessage := startTCPServer(t)
	// init the Sender and write the message
	snd, err := NewTCPSender(host, port, time.Second)
	require.NoError(t, err)

	msg := fmt.Sprint(rand.Int())
	require.NoError(t, snd.Send([]byte(msg)))
	require.NoError(t, snd.Close())

	require.Equal(t, []byte(msg), getMessage())
}

func startTCPServer(t *testing.T) (host, port string, getMessage func() []byte) {
	// spawn a TCP server that would be listening to incoming connections and
	// will save the message in bytesRead
	lis, err := net.Listen("tcp", "localhost:")
	require.NoError(t, err)

	var (
		bytesRead bytes.Buffer
		wg        sync.WaitGroup
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := lis.Accept()
		require.NoError(t, err)

		_, err = io.Copy(&bytesRead, conn)
		require.NoError(t, err)
	}()
	addr := strings.Split(lis.Addr().String(), ":")
	return addr[0], addr[1], func() []byte {
		defer lis.Close()
		wg.Wait()
		return bytesRead.Bytes()
	}
}

func parseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
