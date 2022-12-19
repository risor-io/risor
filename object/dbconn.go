package object

type DatabaseConnection struct {
	Conn interface{}
}

func (c *DatabaseConnection) Type() Type {
	return DB_CONNECTION
}

func (c *DatabaseConnection) Inspect() string {
	return "<DB_CONNECTION>"
}

func (c *DatabaseConnection) Interface() interface{} {
	return c.Conn
}

func (c *DatabaseConnection) Equals(other Object) Object {
	value := other.Type() == DB_CONNECTION && c.Conn == other.(*DatabaseConnection).Conn
	return NewBool(value)
}

func (c *DatabaseConnection) GetAttr(name string) (Object, bool) {
	return nil, false
}
