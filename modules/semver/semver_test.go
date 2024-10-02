//go:build semver
// +build semver

package semver

import (
	"context"
	"testing"

	"github.com/risor-io/risor/object"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	ctx := context.Background()
	result := Parse(ctx, object.NewString("1.2.3-beta.1+build.123"))
	assert.IsType(t, &object.Map{}, result)
	m := result.(*object.Map)
	assert.Equal(t, int64(1), m.Get("major").(*object.Int).Value())
	assert.Equal(t, int64(2), m.Get("minor").(*object.Int).Value())
	assert.Equal(t, int64(3), m.Get("patch").(*object.Int).Value())
	assert.Equal(t, "beta.1", m.Get("pre").(*object.String).Value())
	assert.Equal(t, "build.123", m.Get("build").(*object.String).Value())
}

func TestBuild(t *testing.T) {
	ctx := context.Background()
	result := Build(ctx, object.NewString("1.2.3+build.123"))
	assert.Equal(t, "build.123", result.(*object.String).Value())
}

func TestValidate(t *testing.T) {
	ctx := context.Background()
	result := Validate(ctx, object.NewString("1.2.3"))
	assert.Nil(t, result)

	result = Validate(ctx, object.NewString("invalid"))
	assert.IsType(t, &object.Error{}, result)
}

func TestMajor(t *testing.T) {
	ctx := context.Background()
	result := Major(ctx, object.NewString("1.2.3"))
	assert.Equal(t, int64(1), result.(*object.Int).Value())
}

func TestMinor(t *testing.T) {
	ctx := context.Background()
	result := Minor(ctx, object.NewString("1.2.3"))
	assert.Equal(t, int64(2), result.(*object.Int).Value())
}

func TestPatch(t *testing.T) {
	ctx := context.Background()
	result := Patch(ctx, object.NewString("1.2.3"))
	assert.Equal(t, int64(3), result.(*object.Int).Value())
}

func TestCompare(t *testing.T) {
	ctx := context.Background()
	result := Compare(ctx, object.NewString("1.2.3"), object.NewString("1.2.4"))
	assert.Equal(t, int64(-1), result.(*object.Int).Value())

	result = Compare(ctx, object.NewString("1.2.3"), object.NewString("1.2.3"))
	assert.Equal(t, int64(0), result.(*object.Int).Value())

	result = Compare(ctx, object.NewString("1.2.4"), object.NewString("1.2.3"))
	assert.Equal(t, int64(1), result.(*object.Int).Value())
}

func TestEquals(t *testing.T) {
	ctx := context.Background()
	result := Equals(ctx, object.NewString("1.2.3"), object.NewString("1.2.3"))
	assert.Equal(t, true, result.(*object.Bool).Value())

	result = Equals(ctx, object.NewString("1.2.3"), object.NewString("1.2.4"))
	assert.Equal(t, false, result.(*object.Bool).Value())

	result = Equals(ctx, object.NewString("1.2.3-1"), object.NewString("1.2.3"))
	assert.Equal(t, false, result.(*object.Bool).Value())
}

func TestPre(t *testing.T) {
	ctx := context.Background()
	result := Pre(ctx, object.NewString("1.2.3-beta.1"))
	assert.Equal(t, "beta.1", result.(*object.String).Value())
}
