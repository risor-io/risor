// Test equivalence: "not in" should equal !("element" in container)
// expected value: true
// expected type: bool

// Both expressions should evaluate to the same result
("element" not in {"key": "value"}) == !("element" in {"key": "value"})