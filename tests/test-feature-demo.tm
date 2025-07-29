// Demonstration of the new "not in" feature
// This replaces the awkward !("element" in map) syntax
// expected value: true  
// expected type: bool

// Before: !("element" in {"key": "value"})
// Now: "element" not in {"key": "value"}
"element" not in {"key": "value", "another": "test"}