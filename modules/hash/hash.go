package base64

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/risor-io/risor/object"
)

func Hash(ctx context.Context, args ...object.Object) object.Object {
	nArgs := len(args)
	if nArgs < 1 || nArgs > 2 {
		return object.Errorf("type error: hash() takes 1 or 2 arguments (%d given)", nArgs)
	}
	data, err := object.AsBytes(args[0])
	if err != nil {
		return err
	}
	alg := "sha256"
	if nArgs == 2 {
		var err *object.Error
		alg, err = object.AsString(args[1])
		if err != nil {
			return err
		}
	}
	var h hash.Hash
	// Hash `data` using the algorithm specified by `alg` and return the result as a byte_slice.
	// Support `alg` values: sha256, sha512, sha1, md5
	switch alg {
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	case "sha1":
		h = sha1.New()
	case "md5":
		h = md5.New()
	default:
		return object.Errorf("type error: hash() algorithm must be one of sha256, sha512, sha1, md5")
	}
	h.Write(data)
	return object.NewByteSlice(h.Sum(nil))

}

func Builtins() map[string]object.Object {
	return map[string]object.Object{
		"hash": object.NewBuiltin("hash", Hash),
	}
}
