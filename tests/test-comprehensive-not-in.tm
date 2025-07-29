// Comprehensive test for "not in" operator - addresses original feature request
// expected value: true
// expected type: bool

// Original feature request: "It's kind of awkward to do if !("element" in map) { } 
// when checking to see if a key exists in a map. I think it would feel a bit more 
// natural to have not for maps"

"element" not in {"key1": "value1", "key2": "value2"}