package model

import (
	"fmt"
)

const BinaryTableSuffix = "_binaries"

// Binary represents record of {prefix}_binaries table
type Binary struct {
	tableName string
	ID        int64
	Name      string
	Data      []byte
	Hash      string
	MimeType  string
}

// SetTablePrefix is setting table prefix
func (b *Binary) SetTablePrefix(prefix string) {
	b.tableName = prefix + BinaryTableSuffix
}

// SetTableName sets name of table
func (b *Binary) SetTableName(tableName string) {
	b.tableName = tableName
}

// TableName returns name of table
func (b *Binary) TableName() string {
	return b.tableName
}

// Get is retrieving model from database
func (b *Binary) Get(appID, memberID int64, name string) (bool, error) {
	return isFound(DBConn.Where("app_id = ? AND member_id = ? AND name = ?", appID, memberID, name).Select("id,name,hash").First(b))
}

// Link returns link to binary data
func (b *Binary) Link() string {
	return fmt.Sprintf(`/data/%s/%d/%s/%s`, b.TableName(), b.ID, "data", b.Hash)
}

// GetByID is retrieving model from db by id
func (b *Binary) GetByID(id int64) (bool, error) {
	return isFound(DBConn.Where("id=?", id).First(b))
}
