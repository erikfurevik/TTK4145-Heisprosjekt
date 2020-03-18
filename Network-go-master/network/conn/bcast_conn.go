// +build !windows

//Denne modulen er hentet rett fra utdelt nettverk, så det må kommenteres i readMe.
//Kommentarer herifra og ned skal fjernes, men kommenterer litt slik at det er lett å finne frem senere

//Er vel bare å drite i det med å få til å kjøre nettverket på windows? Programmet skal vel uansett testes over linux systemer
//og er vel ikke akkurat noe krav om å kjøre det på windows?

package conn

import (
	"net"
	"os"
	"syscall"
)

func DialBroadcastUDP(port int) net.PacketConn {
	s, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	syscall.Bind(s, &syscall.SockaddrInet4{Port: port})

	f := os.NewFile(uintptr(s), "")
	conn, _ := net.FilePacketConn(f)
	f.Close()

	return conn
}
