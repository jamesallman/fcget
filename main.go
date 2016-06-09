package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
)

const (
	API_BASE = "https://support.fortinet.com"
	API_PATH = "/RegistrationAPI/FCWS_RegistrationService.svc/REST/REST_GetAssets"
)

var (
	url    = os.Getenv("FCURL")
	tkn    = os.Getenv("FCTKN")
	verify = true
	sn     string
	client *http.Client
)

func init() {
	flag.StringVar(&url, "url", url, "API url")
	flag.StringVar(&tkn, "token", tkn, "API token")
	flag.BoolVar(&verify, "verify", verify, "SSL certificate verification")
}

func main() {
	_ = "breakpoint"

	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("Serial number parameter is required")
	}

	if url == "" {
		url = API_BASE
	}
	url += API_PATH

	if tkn == "" {
		log.Fatal("API token is required")
	}

	sn = flag.Arg(0)

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !verify},
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client = &http.Client{Transport: tr, Jar: jar}
	client.Transport = tr
	client.Jar = jar

	b, err := json.Marshal(map[string]string{"Token": tkn, "Serial_Number": sn})
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
