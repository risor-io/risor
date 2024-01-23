# jmespath

Module `jmespath` provides json filtering and manipulation using the Jmespath expression syntax.

## Functions

### jmespath

```go filename="Function signature"
jmespath(in object, expression string)
```

Returns the filtered object after the expression has been applied.

```go filename="Example"
>>> data := {
  "locations": [
    {"name": "Seattle", "state": "WA"},
    {"name": "New York", "state": "NY"},
    {"name": "Bellevue", "state": "WA"},
    {"name": "Olympia", "state": "WA"},
  ],
} | jmespath("locations[?state == 'WA'].name | sort(@) | {WashingtonCities: join(', ', @)}")
{
    "WashingtonCities": "Bellevue, Olympia, Seattle"
}

>>> print(jmespath(data, "split(WashingtonCities, ',')")[0])
"Bellevue"
```
