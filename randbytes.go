package randbytes

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net"
	"os"
	"strings"
	"time"
)

var hardwareInfo []byte // = sha256(hostname + pid + start_time + (ip+mac) + random)

func init() {
	buf := make([]byte, 0, 1024*4)
	buf = initWriteHostname(buf)
	buf = initWritePid(buf)
	buf = initWriteStartTime(buf)
	buf = initWriteNetInterface(buf)
	buf = initWriteRand(buf)

	sum := sha256.Sum256(buf)
	hardwareInfo = sum[:]
}

func initWriteHostname(buf []byte) []byte {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return append(buf, hostname...)
}

func initWritePid(buf []byte) []byte {
	pid := os.Getpid()
	var tmp [4]byte
	binary.BigEndian.PutUint32(tmp[:], uint32(pid))
	return append(buf, tmp[:]...)
}

func initWriteStartTime(buf []byte) []byte {
	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], uint64(time.Now().UnixNano()))
	return append(buf, tmp[:]...)
}

func initWriteNetInterface(buf []byte) []byte {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	var flags [4]byte

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}

		buf = append(buf, iface.Name...)
		binary.BigEndian.PutUint32(flags[:], uint32(iface.Flags))
		buf = append(buf, flags[:]...)
		buf = append(buf, iface.HardwareAddr...)

		// interface ip addr must contains ip/mask info,
		// so the addr type must be *net.IPNet
		// the struct implements Addr: *IPAddr, *IPNet, *TCPAddr, *UDPAddr, *UnixAddr
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			buf = append(buf, ipnet.IP...)
		}
	}

	return buf
}

func readRand(buf []byte) {
	n := 0
	for i := 0; i < 100 && n < len(buf); i++ {
		m, err := rand.Read(buf[n:])
		if err != nil {
			continue
		}
		n += m
	}
}

func initWriteRand(buf []byte) []byte {
	n := len(buf)
	buf = buf[:cap(buf)]
	readRand(buf[n:])
	return buf
}

func genLocalRandBytes() []byte {
	buf := make([]byte, 1024)
	n := 0
	binary.BigEndian.PutUint64(buf[n:], uint64(time.Now().UnixNano()))
	n += 8
	n += copy(buf[n:], hardwareInfo)
	readRand(buf[n:])
	return buf
}

// New returns a random bytes of size.
func New(size int) []byte {
	if size <= 0 {
		return nil
	}
	p := make([]byte, size)
	for n := 0; n < len(p); {
		s := sha256.Sum256(genLocalRandBytes())
		n += copy(p[n:], s[:])
	}
	return p
}

// UUID returns a uuid string like:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func UUID() string {
	sum := md5.Sum(genLocalRandBytes())
	uuid := hex.EncodeToString(sum[:])
	return strings.Join([]string{
		uuid[0:8],
		uuid[8:12],
		uuid[12:16],
		uuid[16:20],
		uuid[20:32],
	}, "-")
}
