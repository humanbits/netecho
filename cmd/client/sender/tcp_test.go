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

	// init the Sender and write the message
	hostPort := strings.Split(lis.Addr().String(), ":")
	snd, err := NewTCPSender(hostPort[0], parseInt(hostPort[1]), time.Second)
	require.NoError(t, err)

	msg := fmt.Sprint(rand.Int())
	require.NoError(t, snd.Send([]byte(msg)))
	require.NoError(t, snd.Close())

	// wait until the message is ready and assess it
	wg.Wait()
	require.Equal(t, []byte(msg), bytesRead.Bytes())
}

func parseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
