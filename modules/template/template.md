# template

String templating functionality.

## Builtins

### render

```go filename="Function signature"
render(data object, template string) string
```

Returns the rendered template as a string.
It includes all the sprig lib functions.
You can access environment variables from the template under .Env and the passed values will be available under .Values in the template

If compiled with [`-tags k8s`](https://github.com/risor-io/risor#build-and-install-the-cli-from-source),
it also includes a k8sLookup function to get values from k8s objects.

```go filename="Example"
>>> fetch("http://ipinfo.io").json() | render("You are in {{ .Values.city }}, region {{ .Values.region }} in {{ .Values.timezone }}")
"You are in Dublin, region Leinster in Europe/Dublin"
```

## Functions

### new

```go filename="Function signature"
new(name string) template
```

Instanciates a new template object with the given name.

```go filename="Example"
tpl :=  template.new("test")
tpl.delims("{%", "%}")
```

### add

```go filename="Function signature"
add(name string, template string)
```

Adds a named template to the template object

```go filename="Example"
tpl.add("ipinfo", "You are in {% .city %}, region {% .region %} in {% .timezone %}")
```

### execute_template

```go filename="Function signature"
execute_template(data object, name string) string
```

Renders the given named template into a string.

```go filename="Example"
>>> tpl.execute_template(fetch("http://ipinfo.io").json(), "ipinfo")
"You are in Dublin, region Leinster in Europe/Dublin"
```

### parse

```go filename="Function signature"
parse(template string)
```

Parses a template into the template object

```go filename="Example"
tpl.parse("You are in {% .city %}, region {% .region %} in {% .timezone %}")
```

### execute

```go filename="Function signature"
execute(data object) string
```

Renders the templates into a string.

```go filename="Example"
>>> tpl.execute(fetch("http://ipinfo.io").json())
"You are in Dublin, region Leinster in Europe/Dublin"
```
