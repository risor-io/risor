package net

import (
	"context"
	"net"

	"github.com/risor-io/risor/internal/arg"
	"github.com/risor-io/risor/object"
)

func LookupAddr(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.lookup_addr", 1, args); err != nil {
		return err
	}
	addr, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	names, netErr := net.LookupAddr(addr)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return object.NewStringList(names)
}

func LookupCNAME(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.lookup_cname", 1, args); err != nil {
		return err
	}
	addr, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	cname, netErr := net.LookupCNAME(addr)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return object.NewString(cname)
}

func LookupHost(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.lookup_host", 1, args); err != nil {
		return err
	}
	host, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	addrs, netErr := net.LookupHost(host)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return object.NewStringList(addrs)
}

func LookupPort(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.lookup_port", 2, args); err != nil {
		return err
	}
	network, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	service, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	port, netErr := net.LookupPort(network, service)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return object.NewInt(int64(port))
}

func LookupTXT(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.lookup_txt", 1, args); err != nil {
		return err
	}
	domain, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	txts, netErr := net.LookupTXT(domain)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return object.NewStringList(txts)
}

func ParseCIDR(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.parse_cidr", 1, args); err != nil {
		return err
	}
	cidr, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	_, ipNet, netErr := net.ParseCIDR(cidr)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return NewIPNet(ipNet)
}

func SplitHostPort(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.split_host_port", 1, args); err != nil {
		return err
	}
	hostPort, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	host, port, netErr := net.SplitHostPort(hostPort)
	if netErr != nil {
		return object.NewError(netErr)
	}
	return object.NewList([]object.Object{
		object.NewString(host),
		object.NewString(port),
	})
}

func JoinHostPort(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.join_host_port", 2, args); err != nil {
		return err
	}
	host, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	port, err := object.AsString(args[1])
	if err != nil {
		return err
	}
	return object.NewString(net.JoinHostPort(host, port))
}

func InterfaceAddrs(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.interface_addrs", 0, args); err != nil {
		return err
	}
	addrs, netErr := net.InterfaceAddrs()
	if netErr != nil {
		return object.NewError(netErr)
	}
	var addrStrs []object.Object
	for _, addr := range addrs {
		addrStrs = append(addrStrs, object.NewString(addr.String()))
	}
	return object.NewList(addrStrs)
}

func ParseIP(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.parse_ip", 1, args); err != nil {
		return err
	}
	ipStr, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return object.Errorf("invalid IP address: %q", ipStr)
	}
	return NewIP(ip)
}

func LookupIP(ctx context.Context, args ...object.Object) object.Object {
	if err := arg.Require("net.lookup_ip", 1, args); err != nil {
		return err
	}
	host, err := object.AsString(args[0])
	if err != nil {
		return err
	}
	ips, netErr := net.LookupIP(host)
	if netErr != nil {
		return object.NewError(netErr)
	}
	var ipObjs []object.Object
	for _, ip := range ips {
		ipObjs = append(ipObjs, NewIP(ip))
	}
	return object.NewList(ipObjs)
}

func Module() *object.Module {
	return object.NewBuiltinsModule("net", map[string]object.Object{
		"interface_addrs": object.NewBuiltin("interface_addrs", InterfaceAddrs),
		"join_host_port":  object.NewBuiltin("join_host_port", JoinHostPort),
		"lookup_addr":     object.NewBuiltin("lookup_addr", LookupAddr),
		"lookup_cname":    object.NewBuiltin("lookup_cname", LookupCNAME),
		"lookup_host":     object.NewBuiltin("lookup_host", LookupHost),
		"lookup_ip":       object.NewBuiltin("lookup_ip", LookupIP),
		"lookup_port":     object.NewBuiltin("lookup_port", LookupPort),
		"lookup_txt":      object.NewBuiltin("lookup_txt", LookupTXT),
		"parse_cidr":      object.NewBuiltin("parse_cidr", ParseCIDR),
		"parse_ip":        object.NewBuiltin("parse_ip", ParseIP),
		"split_host_port": object.NewBuiltin("split_host_port", SplitHostPort),
	})
}
