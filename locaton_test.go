package WebGuard

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://192.168.248.1:8443/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
