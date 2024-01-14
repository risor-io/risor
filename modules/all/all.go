package all

import (
	"github.com/risor-io/risor/builtins"
	modBase64 "github.com/risor-io/risor/modules/base64"
	modBytes "github.com/risor-io/risor/modules/bytes"
	modExec "github.com/risor-io/risor/modules/exec"
	modFilepath "github.com/risor-io/risor/modules/filepath"
	modFmt "github.com/risor-io/risor/modules/fmt"
	modHash "github.com/risor-io/risor/modules/hash"
	modHttp "github.com/risor-io/risor/modules/http"
	modJson "github.com/risor-io/risor/modules/json"
	modMath "github.com/risor-io/risor/modules/math"
	modOs "github.com/risor-io/risor/modules/os"
	modRand "github.com/risor-io/risor/modules/rand"
	modRegexp "github.com/risor-io/risor/modules/regexp"
	modStrconv "github.com/risor-io/risor/modules/strconv"
	modStrings "github.com/risor-io/risor/modules/strings"
	modThreads "github.com/risor-io/risor/modules/threads"
	modTime "github.com/risor-io/risor/modules/time"
	"github.com/risor-io/risor/object"
)

func Builtins() map[string]object.Object {
	result := map[string]object.Object{
		"base64":   modBase64.Module(),
		"bytes":    modBytes.Module(),
		"exec":     modExec.Module(),
		"filepath": modFilepath.Module(),
		"fmt":      modFmt.Module(),
		"http":     modHttp.Module(),
		"json":     modJson.Module(),
		"math":     modMath.Module(),
		"os":       modOs.Module(),
		"rand":     modRand.Module(),
		"regexp":   modRegexp.Module(),
		"strconv":  modStrconv.Module(),
		"strings":  modStrings.Module(),
		"threads":  modThreads.Module(),
		"time":     modTime.Module(),
	}
	for k, v := range modHttp.Builtins() {
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
