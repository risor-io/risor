package vm

import (
	"context"
	"testing"

	"github.com/risor-io/risor/compiler"
	"github.com/risor-io/risor/object"
)

func TestDynamicStackAllocation(t *testing.T) {
	// Test with dynamic allocation enabled
	vm, err := NewEmpty()
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	// Apply dynamic memory limits
	limits := DefaultMemoryLimits()
	limits.InitialStackSize = 4  // Start very small to test growth
	limits.MaxStackSize = 64

	option := WithMemoryLimits(limits)
	if err := vm.applyOptions([]Option{option}); err != nil {
		t.Fatalf("Failed to apply options: %v", err)
	}

	// Verify dynamic allocation is enabled
	if !vm.useDynamic {
		t.Fatal("Dynamic allocation should be enabled")
	}

	// Verify initial stack capacity
	if vm.stackCap != 4 {
		t.Errorf("Expected initial stack capacity of 4, got %d", vm.stackCap)
	}

	// Test stack growth by pushing more items than initial capacity
	testObjects := []object.Object{
		object.NewInt(1),
		object.NewInt(2),
		object.NewInt(3),
		object.NewInt(4),
		object.NewInt(5),  // This should trigger growth
		object.NewInt(6),
		object.NewInt(7),
		object.NewInt(8),
		object.NewInt(9),  // This should trigger another growth
	}

	for i, obj := range testObjects {
		vm.push(obj)
		if vm.sp != i {
			t.Errorf("After push %d, expected sp=%d, got sp=%d", i, i, vm.sp)
		}
	}

	// Verify stack has grown
	if vm.stackCap <= 4 {
		t.Errorf("Stack should have grown beyond initial capacity of 4, got %d", vm.stackCap)
	}

	// Verify all objects can be popped correctly
	for i := len(testObjects) - 1; i >= 0; i-- {
		obj := vm.pop()
		expected := testObjects[i]
		if !object.Equal(obj, expected) {
			t.Errorf("Pop %d: expected %v, got %v", i, expected, obj)
		}
		if vm.sp != i-1 {
			t.Errorf("After pop %d, expected sp=%d, got sp=%d", i, i-1, vm.sp)
		}
	}
}

func TestLegacyModeCompatibility(t *testing.T) {
	// Test with legacy mode (default)
	vm, err := NewEmpty()
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	// Verify legacy mode is default
	if vm.useDynamic {
		t.Fatal("Legacy mode should be default")
	}

	// Test basic stack operations work in legacy mode
	vm.push(object.NewInt(42))
	if vm.sp != 0 {
		t.Errorf("Expected sp=0, got sp=%d", vm.sp)
	}

	obj := vm.pop()
	if !object.Equal(obj, object.NewInt(42)) {
		t.Errorf("Expected 42, got %v", obj)
	}
	if vm.sp != -1 {
		t.Errorf("Expected sp=-1, got sp=%d", vm.sp)
	}
}

func TestFrameAllocation(t *testing.T) {
	// Test with dynamic frame allocation
	vm, err := NewEmpty()
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	// Apply dynamic memory limits with small initial frame count
	limits := DefaultMemoryLimits()
	limits.InitialFrameCount = 2
	limits.MaxFrameCount = 16

	option := WithMemoryLimits(limits)
	if err := vm.applyOptions([]Option{option}); err != nil {
		t.Fatalf("Failed to apply options: %v", err)
	}

	// Verify frame capacity
	if vm.framesCap != 2 {
		t.Errorf("Expected initial frame capacity of 2, got %d", vm.framesCap)
	}

	// Test frame growth by ensuring capacity
	if err := vm.ensureFrameCapacity(8); err != nil {
		t.Fatalf("Failed to ensure frame capacity: %v", err)
	}

	// Verify frames have grown
	if vm.framesCap < 8 {
		t.Errorf("Frame capacity should be at least 8, got %d", vm.framesCap)
	}
}

func TestMemoryLimits(t *testing.T) {
	vm, err := NewEmpty()
	if err != nil {
		t.Fatalf("Failed to create VM: %v", err)
	}

	// Apply very restrictive limits
	limits := VMMemoryLimits{
		MaxStackSize:      4,
		MaxFrameCount:     2,
		MaxArgsLimit:      3,
		InitialStackSize:  2,
		InitialFrameCount: 1,
	}

	option := WithMemoryLimits(limits)
	if err := vm.applyOptions([]Option{option}); err != nil {
		t.Fatalf("Failed to apply options: %v", err)
	}

	// Test stack limit enforcement
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when exceeding stack limit")
		}
	}()

	// Push beyond the limit (should panic)
	for i := 0; i < 6; i++ {
		vm.push(object.NewInt(int64(i)))
	}
}