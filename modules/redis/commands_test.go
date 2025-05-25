package redis

import (
	"context"
	"github.com/risor-io/risor/object"
	"testing"
)

func TestRedisPing(t *testing.T) {
	obj, ok := client.GetAttr("ping")
	if !ok {
		t.Fatal("expected 'ping' method to exist on Client client")
	}
	fn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'ping' to be a builtin function, got %T", obj)
	}
	result := fn.Call(context.Background())
	res, err := object.AsString(result)
	if err != nil {
		t.Fatalf("expected ping result to be a string, got %T", result)
	}
	if res != "PONG" {
		t.Errorf("expected ping result to be 'PONG', got %s", res)
	}
}

func TestRedisCommands(t *testing.T) {
	key := object.NewString("test_key")
	value := object.NewString("test_value")
	testRedisSet(t, key, value)
	testRedisGet(t, key)
	testRedisExists(t, key)
	testRedisExpire(t, key)
	testRedisTTL(t, key)
	testRedisDel(t, key)
}

func TestRedisKeys(t *testing.T) {
	testRedisSet(t, object.NewString("test_key:1"), object.NewString("test_value:1"))
	testRedisSet(t, object.NewString("test_key:2"), object.NewString("test_value:2"))

	// Test Redis keys operation
	obj, ok := client.GetAttr("keys")
	if !ok {
		t.Fatal("expected 'keys' method to exist on Client client")
	}
	keysFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'keys' to be a builtin function, got %T", obj)
	}
	pattern := object.NewString("*")
	result := keysFn.Call(context.Background(), pattern)
	res, err := object.AsStringSlice(result)
	if err != nil {
		t.Fatalf("expected keys result to be an array, got %T", result)
	}

	if len(res) != 2 {
		t.Errorf("expected keys result to contain 2 elements, got %d", len(res))
	}
}

func TestRedisIncrDecr(t *testing.T) {
	key := object.NewString("counter")

	// Test Redis increment operation
	obj, ok := client.GetAttr("incr")
	if !ok {
		t.Fatal("expected 'incr' method to exist on Client client")
	}
	incrFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'incr' to be a builtin function, got %T", obj)
	}
	incrResult := incrFn.Call(context.Background(), key)
	incrValue, err := object.AsInt(incrResult)
	if err != nil {
		t.Fatalf("expected incr result to be an integer, got %T", incrResult)
	}
	if incrValue != 1 {
		t.Errorf("expected incr result to be 2, got %d", incrValue)
	}

	// Test Redis decrement operation
	obj, ok = client.GetAttr("decr")
	if !ok {
		t.Fatal("expected 'decr' method to exist on Client client")
	}
	decrFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'decr' to be a builtin function, got %T", obj)
	}
	decrResult := decrFn.Call(context.Background(), key)
	decrValue, err := object.AsInt(decrResult)
	if err != nil {
		t.Fatalf("expected decr result to be an integer, got %T", decrResult)
	}
	if decrValue != 0 {
		t.Errorf("expected decr result to be 1, got %d", decrValue)
	}
}

func TestRedisFlushDB(t *testing.T) {
	// Test Redis flushDB operation
	obj, ok := client.GetAttr("flushdb")
	if !ok {
		t.Fatal("expected 'flushdb' method to exist on Client client")
	}
	flushFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'flushdb' to be a builtin function, got %T", obj)
	}
	result := flushFn.Call(context.Background())
	res, err := object.AsString(result)
	if err != nil {
		t.Fatalf("expected flushdb result to be a string, got %T", result)
	}
	if res != "OK" {
		t.Errorf("expected flushdb to succeed, got %v", res)
	}
}

func testRedisSet(t *testing.T, key, value object.Object) {
	// Test Redis set operation
	obj, ok := client.GetAttr("set")
	if !ok {
		t.Fatal("expected 'set' method to exist on Client client")
	}
	setFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'set' to be a builtin function, got %T", obj)
	}
	result := setFn.Call(context.Background(), key, value)
	res, err := object.AsString(result)
	if err != nil {
		t.Fatalf("expected set result to be a string, got %T", result)
	}
	if res != "OK" {
		t.Errorf("expected set to succeed, got %v", res)
	}
}

func testRedisGet(t *testing.T, key object.Object) {
	// Test Redis get operation
	obj, ok := client.GetAttr("get")
	if !ok {
		t.Fatal("expected 'get' method to exist on Client client")
	}
	getFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'get' to be a builtin function, got %T", obj)
	}
	result := getFn.Call(context.Background(), key)
	res, err := object.AsString(result)
	if err != nil {
		t.Fatalf("expected get result to be a string, got %T", result)
	}
	if res != "test_value" {
		t.Errorf("expected get result to be 'test_value', got %s", res)
	}

}

func testRedisExists(t *testing.T, key object.Object) {
	// Test Redis exists operation
	obj, ok := client.GetAttr("exists")
	if !ok {
		t.Fatal("expected 'exists' method to exist on Client client")
	}
	existsFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'exists' to be a builtin function, got %T", obj)
	}
	existsResult := existsFn.Call(context.Background(), key)
	res, err := object.AsInt(existsResult)
	if err != nil {
		t.Fatalf("expected exists result to be a string, got %T", existsResult)
	}
	if res != 1 {
		t.Errorf("expected exists result to be '1', got %d", res)
	}
}

func testRedisExpire(t *testing.T, key object.Object) {
	// Test Redis expire operation
	obj, ok := client.GetAttr("expire")
	if !ok {
		t.Fatal("expected 'expire' method to exist on Client client")
	}
	expireFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'expire' to be a builtin function, got %T", obj)
	}
	result := expireFn.Call(context.Background(), key, object.NewInt(300))
	res, err := object.AsBool(result)
	if err != nil {
		t.Fatalf("expected expire result to be a boolean, got %T", result)
	}
	if !res {
		t.Errorf("expected expire to succeed, got %v", res)
	}
}

func testRedisTTL(t *testing.T, key object.Object) {
	// Test Redis TTL operation
	obj, ok := client.GetAttr("ttl")
	if !ok {
		t.Fatal("expected 'ttl' method to exist on Client client")
	}
	ttlFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'ttl' to be a builtin function, got %T", obj)
	}
	result := ttlFn.Call(context.Background(), key)
	res, err := object.AsInt(result)
	if err != nil {
		t.Fatalf("expected ttl result to be an integer, got %T", result)
	}
	if res <= 0 {
		t.Errorf("expected ttl to be greater than 0, got %d", res)
	}
}

func testRedisDel(t *testing.T, key object.Object) {
	// Test Redis del operation
	obj, ok := client.GetAttr("del")
	if !ok {
		t.Fatal("expected 'del' method to exist on Client client")
	}
	delFn, ok := obj.(*object.Builtin)
	if !ok {
		t.Fatalf("expected 'del' to be a builtin function, got %T", obj)
	}
	result := delFn.Call(context.Background(), key)
	res, err := object.AsInt(result)
	if err != nil {
		t.Fatalf("expected del result to be an integer, got %T", result)
	}
	if res != 1 {
		t.Errorf("expected del to succeed, got %d", res)
	}
}
