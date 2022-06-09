package WebGuard

import (
	"errors"
	"net"
	"regexp"
	"strconv"
	"strings"
)

func IPIntoSlicesB(ip string) ([]string, error) {
	var ip4 = strings.Split(ip, "-")
	var ipA = net.ParseIP(ip4[0])
	if ip4 == nil {
		return []string{}, errors.New("not an ipv4 address")
	}
	var temp []string
	if len(ip4[1]) < 4 {
		iprange, err := strconv.Atoi(ip4[1])
		if ipA == nil || iprange > 255 || err != nil {
			return []string{}, errors.New("input format is not correct")
		}
		var splitip = strings.Split(ip4[0], ".")
		ip1, err1 := strconv.Atoi(splitip[3])
		ip2, err2 := strconv.Atoi(ip4[1])
		prefix := strings.Join(splitip[0:3], ".")
		if ip1 > ip2 || err1 != nil || err2 != nil {
			return []string{}, errors.New("input format is not correct")
		}
		for i := ip1; i <= ip2; i++ {
			temp = append(temp, prefix+"."+strconv.Itoa(i))
		}
	} else {
		var splitip1 = strings.Split(ip4[0], ".")
		var splitip2 = strings.Split(ip4[1], ".")
		if len(splitip1) != 4 || len(splitip2) != 4 {
			return []string{}, errors.New("input format is not correct")
		}
		start, end := [4]int{}, [4]int{}
		for i := 0; i < 4; i++ {
			ip1, err1 := strconv.Atoi(splitip1[i])
			ip2, err2 := strconv.Atoi(splitip2[i])
			if ip1 > ip2 || err1 != nil || err2 != nil {
				return []string{}, errors.New("input format is not correct")
			}
			start[i], end[i] = ip1, ip2
		}
		startNum := start[0]<<24 | start[1]<<16 | start[2]<<8 | start[3]
		endNum := end[0]<<24 | end[1]<<16 | end[2]<<8 | end[3]
		for num := startNum; num <= endNum; num++ {
			ip := strconv.Itoa((num>>24)&0xff) + "." + strconv.Itoa((num>>16)&0xff) + "." + strconv.Itoa((num>>8)&0xff) + "." + strconv.Itoa((num)&0xff)
			temp = append(temp, ip)
		}
	}
	return temp, nil
}

func IPIntoSlicesA(ip string) ([]string, error) {
	var ip4 = net.ParseIP(strings.Split(ip, "/")[0])
	if ip4 == nil {
		return []string{}, errors.New("not an ipv4 address")
	}
	var mark = strings.Split(ip, "/")[1]
	var temp []string
	var err error
	switch mark {
	case "24":
		var ip3 = strings.Join(strings.Split(ip[:], ".")[0:3], ".")
		for i := 0; i <= 255; i++ {
			temp = append(temp, ip3+"."+strconv.Itoa(i))
		}
		err = nil
	case "16":
		var ip2 = strings.Join(strings.Split(ip[:], ".")[0:2], ".")
		for i := 0; i <= 255; i++ {
			for j := 0; j <= 255; j++ {
				temp = append(temp, ip2+"."+strconv.Itoa(i)+"."+strconv.Itoa(j))
			}
		}
		err = nil
	default:
		temp = []string{}
		err = errors.New("not currently supported")
	}
	return temp, err
}

func IPIntoSlices(ip string) ([]string, error) {
	reg := regexp.MustCompile(`[a-zA-Z]+`)
	switch {
	case strings.Count(ip, "/") == 1:
		return IPIntoSlicesA(ip)
	case strings.Count(ip, "-") == 1:
		return IPIntoSlicesB(ip)
	case reg.MatchString(ip):
		_, err := net.LookupHost(ip)
		if err != nil {
			return []string{}, err
		}
		return []string{ip}, nil
	default:
		var isis = net.ParseIP(ip)
		if isis == nil {
			return []string{}, errors.New("input format is not correct")
		}
		return []string{ip}, nil
	}
}
