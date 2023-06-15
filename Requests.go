package FlowX

import (
	"context"
	"encoding/base64"
	"golang.org/x/net/proxy"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ProxyType int

const (
	SOCKS5 ProxyType = iota
	SOCKS4
	HTTP
)

type ProxyConfig struct {
	Type     ProxyType
	Address  string
	Username string
	Password string
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/B08C390",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:54.0) Gecko/20100101 Firefox/54.0",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/17.17134",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko",
}

func AddRequestHeader(request *http.Request, key, value string) {
	request.Header.Add(key, value)
}

func AddRequestHeaders(request *http.Request, headers map[string]string) {
	for key, value := range headers {
		request.Header.Add(key, value)
	}
}

func SetUserAgent(request *http.Request, agent string) {
	request.Header.Set("User-Agent", agent)
}

func RandomUserAgent(request *http.Request) {
	rand.Seed(time.Now().UnixNano())
	agent := userAgents[rand.Intn(len(userAgents))]
	request.Header.Set("User-Agent", agent)
}

func GetResponseBody(response *http.Response) (string, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	return string(body), nil
}

func CreateGetRequest(url string) (*http.Request, error) {
	request, err := http.NewRequest("GET", url, nil)
	return request, err
}

func CreatePostRequest(url string, contentType string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return request, nil
}

func CreateHTTPClient(autoRedirect bool, proxyConfig *ProxyConfig) (*http.Client, error) {
	var transport http.Transport

	if proxyConfig != nil {
		switch proxyConfig.Type {
		case SOCKS5:
			dialer, err := proxy.SOCKS5("tcp", proxyConfig.Address, nil, proxy.Direct)
			if err != nil {
				return nil, err
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		case SOCKS4:
			dialer, err := proxy.SOCKS5("tcp", proxyConfig.Address, nil, &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			})
			if err != nil {
				return nil, err
			}
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		case HTTP:
			proxyURL, err := url.Parse(proxyConfig.Address)
			if err != nil {
				return nil, err
			}
			transport.Proxy = http.ProxyURL(proxyURL)

			if proxyConfig.Username != "" && proxyConfig.Password != "" {
				auth := proxyConfig.Username + ":" + proxyConfig.Password
				basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
				transport.ProxyConnectHeader = http.Header{}
				transport.ProxyConnectHeader.Set("Proxy-Authorization", basicAuth)
			}
		}
	}

	client := &http.Client{Transport: &transport}

	if autoRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client, nil
}

func UseProxy(proxyConfig *ProxyConfig, request *http.Request) error {
	client, err := CreateHTTPClient(false, proxyConfig)
	if err != nil {
		return err
	}

	_, err = ExecuteRequest(client, request)
	return err
}

func ExecuteRequest(client *http.Client, request *http.Request) (*http.Response, error) {
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
