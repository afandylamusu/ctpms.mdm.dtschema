package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/afandylamusu/ctpms.mdm.dtschema/dataset"
	"github.com/afandylamusu/ctpms.mdm.dtschema/models"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	tableName  = "datasets"
)

type mysqlDataSetRepository struct {
	Conn *sql.DB
}

// NewMysqlDataSetRepository will create an object that represent the dataset.Repository interface
func NewMysqlDataSetRepository(Conn *sql.DB) dataset.Repository {
	return &mysqlDataSetRepository{Conn}
}

func (m *mysqlDataSetRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.DataSet, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.DataSet, 0)
	for rows.Next() {
		t := new(models.DataSet)
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlDataSetRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.DataSet, string, error) {
	query := `SELECT id, name, updated_at, created_at
  						FROM ` + tableName + ` WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, err
}
func (m *mysqlDataSetRepository) GetByID(ctx context.Context, id int64) (res *models.DataSet, err error) {
	query := `SELECT id, name, updated_at, created_at
  						FROM ` + tableName + ` WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *mysqlDataSetRepository) GetByTitle(ctx context.Context, title string) (res *models.DataSet, err error) {
	query := `SELECT id, name, updated_at, created_at
  						FROM ` + tableName + ` WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}
	return
}

func (m *mysqlDataSetRepository) Store(ctx context.Context, a *models.DataSet) error {
	query := `INSERT ` + tableName + ` SET name=? , updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Name, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func (m *mysqlDataSetRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM ` + tableName + ` WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {

		return err
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
		return err
	}

	return nil
}
func (m *mysqlDataSetRepository) Update(ctx context.Context, ar *models.DataSet) error {
	query := `UPDATE ` + tableName + ` set name=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, ar.Name, ar.UpdatedAt, ar.ID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
