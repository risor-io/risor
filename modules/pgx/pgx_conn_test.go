package pgx

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/risor-io/risor"
	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
	"os"
	"runtime"
	"testing"
	"text/template"
	"time"
)

const (
	//DEFAULT_PG_HOST     = "127.0.0.1"
	//DEFAULT_PG_PORT     = "5432"
	//DEFAULT_PG_USER     = "risor"
	//DEFAULT_PG_PASSWORD = "risorpw"
	//DEFAULT_PG_DB       = "risordb"
	DEFAULT_PG_HOST            = "192.168.1.244"
	DEFAULT_PG_PORT            = "5433"
	DEFAULT_PG_USER            = "uas"
	DEFAULT_PG_PASSWORD        = "uas"
	DEFAULT_PG_DB              = "uas"
	DEFAULT_PG_CONNECT_TIMEOUT = "3s"
)

var (
	pgHost           = DEFAULT_PG_HOST
	pgPort           = DEFAULT_PG_PORT
	pgUser           = DEFAULT_PG_USER
	pgPassword       = DEFAULT_PG_PASSWORD
	pgDB             = DEFAULT_PG_DB
	pgConnectTimeout = DEFAULT_PG_CONNECT_TIMEOUT
	pgConnStr        = ""
	enabled          = false
)

func envOrDef(env, def string) string {
	evar := os.Getenv(env)
	if evar == "" {
		return def
	}
	return evar
}

func init() {
	pgHost = envOrDef("PG_HOST", DEFAULT_PG_HOST)
	pgPort = envOrDef("PG_PORT", DEFAULT_PG_PORT)
	pgUser = envOrDef("PG_USER", DEFAULT_PG_USER)
	pgPassword = envOrDef("PG_PASSWORD", DEFAULT_PG_PASSWORD)
	pgDB = envOrDef("PG_DB", DEFAULT_PG_DB)
	pgConnectTimeout = envOrDef("PG_CONNECT_TIMEOUT", DEFAULT_PG_CONNECT_TIMEOUT)
	// postgres://{user}:{pass}@{host}:{port}/{db}
	pgConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", pgUser, pgPassword, pgHost, pgPort, pgDB)
	pgx.Connect(context.Background(), pgConnStr)
	if connTimeout, err := time.ParseDuration(pgConnectTimeout); err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
		defer cancel()
		if conn, err := pgx.Connect(ctx, pgConnStr); err == nil {
			defer conn.Close(context.Background())
			enabled = true
		}
	} else {
		fmt.Printf("Invalid connection timeout: %s\n", pgConnectTimeout)
	}
}

func currentTest() string {
	// Get the current function name using runtime.FuncForPC
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

func compileTemplate(name, source string, args map[string]string) string {
	if args == nil {
		return source
	}
	if templ, err := template.New(name).Parse(source); err != nil {
		panic(err)
	} else {
		var b bytes.Buffer
		if err := templ.Execute(&b, args); err != nil {
			panic(err)
		}
		return fmt.Sprintf(`p := pgx.connect("%s")%s
		p.close()`, pgConnStr, b.String())
	}
}

func Exec(script string, args map[string]string) (object.Object, error) {
	testName := currentTest()
	startTime := time.Now()
	var source string
	if args == nil {
		source = fmt.Sprintf(`p := pgx.connect("%s")%s
		p.close()`, pgConnStr, script)
	} else {
		source = compileTemplate(testName, script, args)
	}
	fmt.Printf("SOURCE:\n===================================\n%s\n===================================\n", source)
	ctx := context.Background()
	defer func() {
		fmt.Printf("Execution: test=%s, elapsed=%s\n", testName, time.Since(startTime))
	}()
	return risor.Eval(ctx, source, risor.WithGlobals(map[string]any{
		"pgx": Module(),
	}))
}

func TestPgx_BasicFunc(t *testing.T) {
	if !enabled {
		t.Skip("No database connected")
		return
	}
	result, err := Exec(BasicFunc, nil)
	require.Nil(t, err)
	require.NotNil(t, result)
}

func TestPgx_GenSeries100000Func(t *testing.T) {
	if !enabled {
		t.Skip("No database connected")
		return
	}
	result, err := Exec(GenSeriesFunc100000, nil)
	require.Nil(t, err)
	require.NotNil(t, result)
}

func TestPgx_GenSeriesParamFunc(t *testing.T) {
	if !enabled {
		t.Skip("No database connected")
		return
	}
	start := 1
	end := 100000
	expected := end - start + 1
	result, err := Exec(GenSeriesFuncParam, map[string]string{
		"start":    fmt.Sprintf("%d", start),
		"end":      fmt.Sprintf("%d", end),
		"expected": fmt.Sprintf("%d", expected),
	})
	require.Nil(t, err)
	require.NotNil(t, result)
}

// Same as TestPgx_GenSeriesParamFunc but using the query() method and counting the rows returned
func TestPgx_GenSeriesParam(t *testing.T) {
	if !enabled {
		t.Skip("No database connected")
		return
	}
	start := 1
	end := 100000
	expected := end - start + 1
	result, err := Exec(GenSeriesParam, map[string]string{
		"start":    fmt.Sprintf("%d", start),
		"end":      fmt.Sprintf("%d", end),
		"expected": fmt.Sprintf("%d", expected),
	})
	require.Nil(t, err)
	require.NotNil(t, result)
}

const (
	BasicFunc = `
	cnt := 0
    rcnt := p.queryFunc("SELECT * FROM pg_stat_activity", func(row) {
		printf("row: %v\n", row)
		cnt++
		return true
	})
	printf("rows returned: %d\n", cnt)
	assert(rcnt == cnt, "row count mismatch")`

	GenSeriesFunc100000 = `
	cnt := 0
    rcnt := p.queryFunc("SELECT * FROM generate_series(1, 100000)", func(row) {
		cnt++
		return true
	})
	printf("rows returned: %d\n", cnt)
	assert(rcnt == cnt, "row count mismatch")`

	GenSeriesFuncParam = `
	cnt := 0
    rcnt := p.queryFunc("SELECT * FROM generate_series($1::integer, $2::integer)", func(row) {
		cnt++
		return true
	}, {{.start}}, {{.end}})
	printf("rows returned: %d\n", cnt)
	assert(rcnt == cnt, "row count mismatch")
	assert(rcnt == {{.expected}}, "row count mismatch to expected")`

	GenSeriesParam = `
	cnt := 0
    rows := p.query("SELECT * FROM generate_series($1::integer, $2::integer)", {{.start}}, {{.end}})
	rows.each(func(row) {
		cnt++
	})
	printf("rows returned: %d\n", cnt)
	assert(cnt == {{.expected}}, "row count mismatch to expected")`
)
