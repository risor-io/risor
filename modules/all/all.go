package all

import (
	"github.com/cloudcmds/tamarin/v2/builtins"
	modBase64 "github.com/cloudcmds/tamarin/v2/modules/base64"
	modBytes "github.com/cloudcmds/tamarin/v2/modules/bytes"
	modFetch "github.com/cloudcmds/tamarin/v2/modules/fetch"
	modFmt "github.com/cloudcmds/tamarin/v2/modules/fmt"
	modHash "github.com/cloudcmds/tamarin/v2/modules/hash"
	modImage "github.com/cloudcmds/tamarin/v2/modules/image"
	modJson "github.com/cloudcmds/tamarin/v2/modules/json"
	modMath "github.com/cloudcmds/tamarin/v2/modules/math"
	modOs "github.com/cloudcmds/tamarin/v2/modules/os"
	modPgx "github.com/cloudcmds/tamarin/v2/modules/pgx"
	modRand "github.com/cloudcmds/tamarin/v2/modules/rand"
	modStrconv "github.com/cloudcmds/tamarin/v2/modules/strconv"
	modStrings "github.com/cloudcmds/tamarin/v2/modules/strings"
	modTime "github.com/cloudcmds/tamarin/v2/modules/time"
	modUuid "github.com/cloudcmds/tamarin/v2/modules/uuid"
	"github.com/cloudcmds/tamarin/v2/object"
)

func Builtins() map[string]object.Object {
	result := map[string]object.Object{
		"math":    modMath.Module(),
		"json":    modJson.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"pgx":     modPgx.Module(),
		"uuid":    modUuid.Module(),
		"os":      modOs.Module(),
		"bytes":   modBytes.Module(),
		"base64":  modBase64.Module(),
		"fmt":     modFmt.Module(),
		"image":   modImage.Module(),
	}
	for k, v := range modFetch.Builtins() {
		result[k] = v
	}
	for k, v := range modFmt.Builtins() {
		result[k] = v
	}
	for k, v := range builtins.Builtins() {
		result[k] = v
	}
	for k, v := range modHash.Builtins() {
		result[k] = v
	}
	return result
}
