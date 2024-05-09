package errctx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithMeta(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"key": "value"}, data)
}

func TestWithMetaAdditional(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")

	err := Wrap(errors.New("a problem"), ctx, "additional", "value")
	data := Unwrap(err)

	assert.Equal(t, map[string]any{
		"key":        "value",
		"additional": "value",
	}, data)
}

func TestWithMetaOverwrite(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")
	ctx = WithMeta(ctx, "key", "value2")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"key": "value2"}, data)
}

func TestWithMetaNested(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value")
	ctx = WithMeta(ctx, "key", "value2")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = WithMeta(ctx, "key", "value3")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"key": "value3"}, data)
}

func TestWithMetaNestedManyKeys(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = WithMeta(ctx, "key2", "value2")
	ctx = WithMeta(ctx, "key3", "value3", "key4", "value4")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
	}, data)
}

func TestWithMetaNestedManyKeysPlusExtraWrappedKV(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "some other", "stuff")
	ctx = WithMeta(ctx, "key2", "value2")
	ctx = WithMeta(ctx, "key3", "value3", "key4", "value4")

	err := Wrap(errors.New("a problem"), ctx, "extra1", "extravalue1", "extra2", "extravalue2")
	data := Unwrap(err)

	assert.Equal(t, map[string]any{
		"key1":   "value1",
		"key2":   "value2",
		"key3":   "value3",
		"key4":   "value4",
		"extra1": "extravalue1",
		"extra2": "extravalue2",
	}, data)
}

func TestWithMetaOddNumberKV(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value", "missingkey")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"key": "value", "!BADKEY": "missingkey"}, data)
}

func TestWithMetaOddNumberWrapKV(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value", "missingkey")

	err := Wrap(errors.New("a problem"), ctx, "wrapkey", "wrapvalue", "missingkey")
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"key": "value", "wrapkey": "wrapvalue", "!BADKEY": "missingkey"}, data)
}

func TestWithMetaOddNumberKVNotString(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "key", "value", 42)

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"key": "value", "!BADKEY": 42}, data)
}

func TestWithMetaOneValueKV(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "missingkey")

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"!BADKEY": "missingkey"}, data)
}

func TestWithMetaOneValueWrapKV(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "missingkey")

	err := Wrap(errors.New("a problem"), ctx, "wrapkey", "wrapvalue", "missingkey")
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"wrapkey": "wrapvalue", "!BADKEY": "missingkey"}, data)
}

func TestWithMetaOneValueEmptyWrapKV(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, "missingkey")

	err := Wrap(errors.New("a problem"), ctx, "missingkey")
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"!BADKEY": "missingkey"}, data)
}

func TestWithMetaOneValueKVNotString(t *testing.T) {
	ctx := context.Background()
	ctx = WithMeta(ctx, 42)

	err := Wrap(errors.New("a problem"), ctx)
	data := Unwrap(err)

	assert.Equal(t, map[string]any{"!BADKEY": 42}, data)
}

func TestWithMetaEmpty(t *testing.T) {
	err := errors.New("a problem")
	data := Unwrap(err)

	assert.Nil(t, data)
}

func TestWithMetaNilContext(t *testing.T) {
	ctx := WithMeta(nil, "key", "value")

	assert.Nil(t, ctx)
}

func TestWithMetaDifferentMapAddress(t *testing.T) {
	ctx := context.Background()
	err := errors.New("a problem")

	ctx1 := WithMeta(ctx, "key1", "value1")
	err1 := Wrap(err, ctx1)

	ctx2 := WithMeta(ctx1, "key2", "value2")
	err2 := Wrap(err1, ctx2)

	ctx3 := WithMeta(ctx2, "key3", 3)
	err3 := Wrap(err1, ctx3)

	data1 := Unwrap(err1)
	data2 := Unwrap(err2)
	data3 := Unwrap(err3)

	assert.Equal(t,
		map[string]any{
			"key1": "value1",
		},
		data1,
		"The first unwrap result is only the key-value pair from the first wrap.",
	)

	assert.Equal(t,
		map[string]any{
			"key1": "value1",
			"key2": "value2",
		},
		data2,
		"The second unwrap result contains all the data merged together.",
	)

	assert.Equal(t,
		map[string]any{
			"key1": "value1",
			"key2": "value2",
			"key3": 3,
		},
		data3,
		"The third unwrap result contains all the data merged together. (different value type)",
	)
}
