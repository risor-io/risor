package dns

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/risor-io/risor/object"
)

func NSLookup(ctx context.Context, args ...object.Object) object.Object {
	numArgs := len(args)
	if numArgs < 1 || numArgs > 3 {
		return object.NewArgsRangeError("nslookup", 1, 3, numArgs)
	}

	addr, argErr := object.AsString(args[0])
	if argErr != nil {
		return argErr
	}

	queryType := "HOST"
	if numArgs > 1 {
		queryType, argErr = object.AsString(args[1])
		if argErr != nil {
			return argErr
		}
	}

	var resolverAddr string
	if numArgs > 2 {
		resolverAddr, argErr = object.AsString(args[2])
		if argErr != nil {
			return argErr
		}
	}

	resolver := net.DefaultResolver
	if resolverAddr != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: 30 * time.Second,
				}
				return d.DialContext(ctx, network, resolverAddr)
			},
		}
	}

	var addrs []string
	var err error

	switch strings.ToUpper(queryType) {
	case "TXT":
		addrs, err = resolver.LookupTXT(ctx, addr)
	case "PTR":
		addrs, err = resolver.LookupAddr(ctx, addr)
	case "CNAME":
		var cname string
		cname, err = resolver.LookupCNAME(ctx, addr)
		addrs = append(addrs, cname)
	case "SRV":
		var res []*net.SRV
		var cname string
		_, res, err = resolver.LookupSRV(ctx, addr, "", "")
		addrs = append(addrs, cname)
		for _, a := range res {
			addrs = append(addrs, a.Target)
		}
	default:
		addrs, err = resolver.LookupHost(ctx, addr)
	}

	if err != nil {
		return object.NewError(err)
	}

	return object.NewStringList(addrs)
}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"nslookup": object.NewBuiltin("nslookup", NSLookup),
	}
}
