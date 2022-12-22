package object

type DatabaseConnection struct {
	conn interface{}
}

func (c *DatabaseConnection) Type() Type {
	return DB_CONNECTION
}

func (c *DatabaseConnection) Inspect() string {
	return "db_connection()"
}

func (c *DatabaseConnection) Interface() interface{} {
	return c.conn
}

func (c *DatabaseConnection) Equals(other Object) Object {
	value := other.Type() == DB_CONNECTION && c.conn == other.(*DatabaseConnection).conn
	return NewBool(value)
}

func (c *DatabaseConnection) GetAttr(name string) (Object, bool) {
	return nil, false
}

func NewDatabaseConnection(conn interface{}) *DatabaseConnection {
	return &DatabaseConnection{conn: conn}
}
