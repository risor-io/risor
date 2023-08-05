package all

import (
	"github.com/risor-io/risor/builtins"
	modAws "github.com/risor-io/risor/modules/aws"
	modBase64 "github.com/risor-io/risor/modules/base64"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modFetch "github.com/risor-io/risor/modules/fetch"
	modFmt "github.com/risor-io/risor/modules/fmt"
	modGoogle "github.com/risor-io/risor/modules/google"
	modHash "github.com/risor-io/risor/modules/hash"
	modImage "github.com/risor-io/risor/modules/image"
	modJson "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modOs "github.com/risor-io/risor/modules/os"
	modPgx "github.com/risor-io/risor/modules/pgx"
	modRand "github.com/risor-io/risor/modules/rand"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modTime "github.com/risor-io/risor/modules/time"
	modUuid "github.com/risor-io/risor/modules/uuid"
	"github.com/risor-io/risor/object"
)

func Builtins() map[string]object.Object {
	result := map[string]object.Object{
		"aws":     modAws.Module(),
		"base64":  modBase64.Module(),
		"bytes":   modBytes.Module(),
		"fmt":     modFmt.Module(),
		"google":  modGoogle.Module(),
		"image":   modImage.Module(),
		"json":    modJson.Module(),
		"math":    modMath.Module(),
		"os":      modOs.Module(),
		"pgx":     modPgx.Module(),
		"rand":    modRand.Module(),
		"strconv": modStrconv.Module(),
		"strings": modStrings.Module(),
		"time":    modTime.Module(),
		"uuid":    modUuid.Module(),
	}

	if awsMod := modAws.Module(); awsMod != nil {
		result["aws"] = awsMod
	}
	if googleMod := modGoogle.Module(); googleMod != nil {
		result["google"] = googleMod
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
	for k, v := range modOs.Builtins() {
		result[k] = v
	}
	return result
}
