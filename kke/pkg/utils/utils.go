package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func checkSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

func Ping(ip string, ports ...string) bool {
	type ICMP struct {
		Type        uint8
		Code        uint8
		Checksum    uint16
		Identifier  uint16
		SequenceNum uint16
	}

	icmp := ICMP{
		Type: 8,
	}

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer

	_ = binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = checkSum(buffer.Bytes())
	buffer.Reset()
	_ = binary.Write(&buffer, binary.BigEndian, icmp)

	Time, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("ip4:icmp", ip, Time)
	if err != nil {
		return exec.Command("ping", ip, "-c", "2", "-i", "1", "-W", "3").Run() == nil
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return false
	}
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	num, err := conn.Read(recvBuf)
	if err != nil {
		return false
	}

	_ = conn.SetReadDeadline(time.Time{})

	return string(recvBuf[0:num]) != ""

}

func PortScan(ip string, ports ...string) (string, bool) {
	res := true
	m := make(map[string]bool)
	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), time.Second*3)
		if err != nil {
			res = false
			m[port] = false
		}
		if conn != nil {
			m[port] = true
			defer conn.Close()
		}
	}
	out, _ := json.Marshal(m)
	return string(out), res
}

func ParseHost(host string) (string, string, string, int, error) {
	host = strings.TrimSpace(host)

	user := "root"
	cred := ""
	addr := "127.0.0.1"
	port := 22
	var err error

	uc := "" // user/cred
	ap := "" // addr:port
	ucap := strings.Split(host, "@")
	if len(ucap) < 2 {
		ap = ucap[0]
	} else {
		uc = strings.Join(ucap[:len(ucap)-1], "@")
		ap = ucap[len(ucap)-1]
	}

	tmp := strings.Split(uc, "/")
	user = tmp[0]
	if len(tmp) >= 2 {
		cred = tmp[1]
	}

	tmp = strings.Split(ap, ":")
	addr = tmp[0]
	if len(tmp) >= 2 {
		port, err = strconv.Atoi(tmp[1])
		if err != nil {
			return user, cred, addr, port, err
		}
	}

	return user, cred, addr, port, nil
}

func StringSplit(str string) []string {
	var out []string
	if str != "" {
		strs := strings.Split(str, " ")
		if len(strs) == 1 {
			strs = strings.Split(str, ",")
		}
		for _, s := range strs {
			out = append(out, strings.TrimSpace(s))
		}
	}

	return out
}
