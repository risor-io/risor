package filepath

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/require"
)

func TestAbs(t *testing.T) {
	ctx := context.Background()
	wd, err := os.Getwd()
	require.NoError(t, err)
	abs := Abs(ctx, object.NewString("foo"))
	require.IsType(t, &object.String{}, abs)
	require.Equal(t, filepath.Join(wd, "foo"), abs.(*object.String).Value())
}

func TestBase(t *testing.T) {
	ctx := context.Background()
	base := Base(ctx, object.NewString("/foo/bar.txt"))
	require.IsType(t, &object.String{}, base)
	require.Equal(t, "bar.txt", base.(*object.String).Value())
}

func TestClean(t *testing.T) {
	ctx := context.Background()
	clean := Clean(ctx, object.NewString("/foo/../foo/bar//baz"))
	require.IsType(t, &object.String{}, clean)
	require.Equal(t, "/foo/bar/baz", clean.(*object.String).Value())
}

func TestDir(t *testing.T) {
	ctx := context.Background()
	dir := Dir(ctx, object.NewString("/foo/bar/baz.txt"))
	require.IsType(t, &object.String{}, dir)
	require.Equal(t, "/foo/bar", dir.(*object.String).Value())
}

func TestExt(t *testing.T) {
	ctx := context.Background()
	ext := Ext(ctx, object.NewString("bar/baz.txt"))
	require.IsType(t, &object.String{}, ext)
	require.Equal(t, ".txt", ext.(*object.String).Value())
}

func TestIsAbs(t *testing.T) {
	ctx := context.Background()
	isAbsTrue := IsAbs(ctx, object.NewString("/foo/bar"))
	require.IsType(t, &object.Bool{}, isAbsTrue)
	require.True(t, isAbsTrue.(*object.Bool).Value())

	isAbsFalse := IsAbs(ctx, object.NewString("foo/bar"))
	require.IsType(t, &object.Bool{}, isAbsFalse)
	require.False(t, isAbsFalse.(*object.Bool).Value())
}

func TestJoin(t *testing.T) {
	ctx := context.Background()
	join := Join(ctx, object.NewString("foo"), object.NewString("bar"), object.NewString("baz.txt"))
	require.IsType(t, &object.String{}, join)
	require.Equal(t, "foo/bar/baz.txt", join.(*object.String).Value())
}

func TestMatch(t *testing.T) {
	ctx := context.Background()
	result := Match(ctx, object.NewString("*.txt"), object.NewString("file.txt"))
	require.IsType(t, &object.Bool{}, result)
	require.True(t, result.(*object.Bool).Value())

	result = Match(ctx, object.NewString("*.txt"), object.NewString("file.jpg"))
	require.IsType(t, &object.Bool{}, result)
	require.False(t, result.(*object.Bool).Value())
}

func TestRel(t *testing.T) {
	ctx := context.Background()
	result := Rel(ctx, object.NewString("/foo"), object.NewString("/foo/bar/baz"))
	require.IsType(t, &object.String{}, result)
	require.Equal(t, "bar/baz", result.(*object.String).Value())
}

func TestSplit(t *testing.T) {
	ctx := context.Background()
	split := Split(ctx, object.NewString("/foo/bar/baz.txt"))
	require.IsType(t, &object.List{}, split)
	l := split.(*object.List)
	items := l.Value()
	require.Len(t, items, 2)
	require.Equal(t, "/foo/bar/", items[0].(*object.String).Value())
	require.Equal(t, "baz.txt", items[1].(*object.String).Value())
}

func TestSplitList(t *testing.T) {
	ctx := context.Background()
	splitList := SplitList(ctx, object.NewString("/foo:/bar:/baz"))
	require.IsType(t, &object.List{}, splitList)
	items := splitList.(*object.List).Value()
	require.Len(t, items, 3)
	require.Equal(t, "/foo", items[0].(*object.String).Value())
	require.Equal(t, "/bar", items[1].(*object.String).Value())
	require.Equal(t, "/baz", items[2].(*object.String).Value())
}

func TestWalkDir(t *testing.T) {

	callFunc := func(ctx context.Context, fn *object.Function, args []object.Object) (object.Object, error) {
		require.FailNow(t, "callFunc should not be called")
		return nil, nil
	}
	ctx := context.Background()
	ctx = object.WithCallFunc(ctx, callFunc)

	var items []string
	result := WalkDir(ctx, object.NewString("testdir"),
		object.NewBuiltin("test", func(ctx context.Context, args ...object.Object) object.Object {
			require.Len(t, args, 3)
			require.IsType(t, &object.String{}, args[0])
			items = append(items, args[0].(*object.String).Value())
			return nil
		}))

	require.Equal(t, object.Nil, result)
	require.Equal(t, []string{
		"testdir",
		"testdir/a",
		"testdir/a/a.txt",
		"testdir/b",
		"testdir/b/b.txt",
	}, items)

	var goldenItems []string
	filepath.WalkDir("testdir", func(path string, info os.DirEntry, err error) error {
		goldenItems = append(goldenItems, path)
		return nil
	})
	require.Equal(t, goldenItems, items)
}
