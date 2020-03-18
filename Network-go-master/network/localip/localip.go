//Denne modulen er hentet rett fra utdelt nettverk, så det må kommenteres i readMe.
//Kommentarer skal fjernes, men kommenterer litt slik at det er lett å finne frem senere
package localip

import (
	"net"
	"strings"
)

var localIP string

//Denne funksjonen vil kun returnere en IP om datamaskinene er online.
func LocalIP() (string, error) {
	if localIP == "" {
		conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{IP: []byte{8, 8, 8, 8}, Port: 53})
		if err != nil {
			return "", err
		}
		defer conn.Close()
		localIP = strings.Split(conn.LocalAddr().String(), ":")[0]
	}
	return localIP, nil
}
