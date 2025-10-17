package database

// sudo -iu postgres
// psql -U eu maindb

// TODO
// - showTable fix sql injection
// - book wrapper struct

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"errors"

	"database/sql"
	_ "github.com/lib/pq"

	"github.com/google/uuid"
)

type Book struct {
	Id uuid.UUID
	Title string
	Author string
	Isbn string
	Price float64 `json:",string"`
	Owner string
}

func NewBook(params ...any) (Book, error) {
	var b Book

	v := reflect.ValueOf(&b).Elem()
	t := v.Type() 
	n := v.NumField()

	if len(params) >= n {
		return Book{}, errors.New("too many parameters")
	}

	for i, p := range params {
		idx := i + 1

		field := v.Field(idx)
		fieldType := t.Field(idx).Type

		if reflect.TypeOf(p) != fieldType {
			return Book{}, fmt.Errorf("type mismatch at parameter %d (got %v instead of %v)", 
				idx, reflect.TypeOf(p), fieldType)
		}

		if field.CanSet() {
			val := reflect.ValueOf(p)
			field.Set(val)
		}
	}

	return b, nil 
}

type Connection struct {
	Db *sql.DB
	cfg Config
}

func Connect(cfg Config) (*Connection, error) {
	db, err := sql.Open("postgres", cfg.ToString())
	if err != nil {
		return &Connection{}, err
	}

	err = db.Ping()
	if err != nil {
		return &Connection{}, err
	}

	fmt.Println("connected to:", cfg.Dbname)

	return &Connection{Db: db, cfg: cfg}, nil
}

func (c *Connection) Close() error {
	return c.Db.Close()
}

func (c *Connection) AddBook(b *Book) error {
    q := `INSERT INTO books (title, author, isbn, price, owner) 
          VALUES ($1, $2, $3, $4, $5)
          RETURNING id`
    
	var id uuid.UUID
    err := c.Db.QueryRow(q, b.Title, b.Author, b.Isbn, b.Price, b.Owner).Scan(&id)
	if err != nil {
		return err
	}
	b.Id = id

    return nil
}

func (c *Connection) GetTables() ([]string, error) {
	q := `
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        ORDER BY table_name
    `

	rows, err := c.Db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tableNames := []string{}
	for rows.Next() {
		var tb string

		err := rows.Scan(&tb)
		if err != nil {
			return nil, err
		}

		tableNames = append(tableNames, tb)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tableNames, nil
}

func GetRow[T any](c *Connection, id uuid.UUID, table string) (T, error) {
	if ts, _ := c.GetTables(); !slices.Contains(ts, table) {
		log.Fatalf("%s is not a table", table)
	}

	var empty T
	q := fmt.Sprintf("SELECT * FROM %s WHERE id='%s'", table, id.String())

	row := c.Db.QueryRow(q)

	var res T
	newItem := reflect.New(reflect.TypeOf(res)).Elem()
	
	var fieldPtrs []any
	for i := 0; i < newItem.NumField(); i++ {
		fieldPtrs = append(fieldPtrs, newItem.Field(i).Addr().Interface())
	}
	
	err := row.Scan(fieldPtrs...)
	if err != nil {
		return empty, err
	}
	
	res = newItem.Interface().(T)

	if err := row.Err(); err != nil {
		return empty, err
	}

	return res, nil
}
