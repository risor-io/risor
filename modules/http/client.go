package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/risor-io/risor/object"
)

func NewHTTPClientFromParams(params *object.Map) (*http.Client, error) {
	client := &http.Client{}

	if storeCookiesObj := params.GetWithDefault("storeCookies", nil); storeCookiesObj != nil {
		storeCookies, errObj := object.AsBool(storeCookiesObj)
		if errObj != nil {
			return nil, errObj.Value()
		}
		if storeCookies {
			jar, err := cookiejar.New(nil)
			if err != nil {
				return nil, err
			}
			client.Jar = jar
		}
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if proxyObj := params.GetWithDefault("proxy", nil); proxyObj != nil {
		proxy, errObj := object.AsString(proxyObj)
		if errObj != nil {
			return nil, errObj.Value()
		}
		if proxy != "" {
			p, err := url.Parse(proxy)
			if err != nil {
				return nil, fmt.Errorf("invalid proxy: %w", err)
			}
			transport.Proxy = http.ProxyURL(p)
		}
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if resolverObj := params.GetWithDefault("resolver", nil); resolverObj != nil {
		resolver, errObj := object.AsString(resolverObj)
		if errObj != nil {
			return nil, errObj.Value()
		}
		if resolver != "" {
			dialer.Resolver = &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					d := net.Dialer{
						Timeout: 30 * time.Second,
					}
					return d.DialContext(ctx, network, resolver)
				},
			}
		}
	}

	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, network, addr)
	}

	client.Transport = transport

	return client, nil
}
