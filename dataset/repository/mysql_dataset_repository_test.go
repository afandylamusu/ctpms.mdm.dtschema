package repository

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

// 	datasetRepo "github.com/afandylamusu/ctpms.mdm.dtschema/dataset/repository"
// 	"github.com/afandylamusu/ctpms.mdm.dtschema/models"
// )

// const (
// 	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
// 	tableName  = "datasets"
// )

// func TestFetch(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer func() {
// 		err = db.Close()
// 		require.NoError(t, err)
// 	}()

// 	mockDataSets := []models.DataSet{
// 		models.DataSet{
// 			ID: 1, Name: "title 1",
// 			UpdatedAt: time.Now(), CreatedAt: time.Now(),
// 		},
// 		models.DataSet{
// 			ID: 2, Name: "title 2",
// 			UpdatedAt: time.Now(), CreatedAt: time.Now(),
// 		},
// 	}

// 	rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
// 		AddRow(mockDataSets[0].ID, mockDataSets[0], t.Name,
// 			mockDataSets[0].UpdatedAt, mockDataSets[0].CreatedAt).
// 		AddRow(mockDataSets[1].ID, mockDataSets[1], t.Name,
// 			mockDataSets[1].UpdatedAt, mockDataSets[1].CreatedAt)

// 	query := "SELECT id, name, updated_at, created_at FROM " + tableName + " WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	a := datasetRepo.NewMysqlDataSetRepository(db)
// 	cursor := datasetRepo.EncodeCursor(mockDataSets[1].CreatedAt)
// 	num := int64(2)
// 	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
// 	assert.NotEmpty(t, nextCursor)
// 	assert.NoError(t, err)
// 	assert.Len(t, list, 2)
// }

// func TestGetByID(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}

// 	defer func() {
// 		err = db.Close()
// 		require.NoError(t, err)
// 	}()

// 	rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
// 		AddRow(1, "name 1", time.Now(), time.Now())

// 	query := "SELECT id,name, updated_at, created_at FROM " + tableName + " WHERE ID = \\?"

// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	a := datasetRepo.NewMysqlDataSetRepository(db)

// 	num := int64(5)
// 	anDataSet, err := a.GetByID(context.TODO(), num)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, anDataSet)
// }

// func TestStore(t *testing.T) {
// 	now := time.Now()
// 	dset := &models.DataSet{
// 		Name:      "Name 1",
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer func() {
// 		err = db.Close()
// 		require.NoError(t, err)
// 	}()

// 	query := "INSERT " + tableName + " SET name=\\? updated_at=\\? , created_at=\\?"
// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(dset.Name, dset.UpdatedAt, dset.CreatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

// 	a := datasetRepo.NewMysqlDataSetRepository(db)

// 	err = a.Store(context.TODO(), dset)
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(12), dset.ID)
// }

// func TestGetByTitle(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer func() {
// 		err = db.Close()
// 		require.NoError(t, err)
// 	}()
// 	rows := sqlmock.NewRows([]string{"id", "name", "updated_at", "created_at"}).
// 		AddRow(1, "name 1", time.Now(), time.Now())

// 	query := "SELECT id, name, updated_at, created_at FROM " + tableName + " WHERE name = \\?"

// 	mock.ExpectQuery(query).WillReturnRows(rows)
// 	a := datasetRepo.NewMysqlDataSetRepository(db)

// 	name := "title 1"
// 	anDataSet, err := a.GetByTitle(context.TODO(), name)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, anDataSet)
// }

// func TestDelete(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer func() {
// 		err = db.Close()
// 		require.NoError(t, err)
// 	}()

// 	query := "DELETE FROM " + tableName + " WHERE id = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

// 	a := datasetRepo.NewMysqlDataSetRepository(db)

// 	num := int64(12)
// 	err = a.Delete(context.TODO(), num)
// 	assert.NoError(t, err)
// }

// func TestUpdate(t *testing.T) {
// 	now := time.Now()
// 	dset := &models.DataSet{
// 		ID:        12,
// 		Name:      "Judul",
// 		CreatedAt: now,
// 		UpdatedAt: now,
// 	}

// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer func() {
// 		err = db.Close()
// 		require.NoError(t, err)
// 	}()

// 	query := "UPDATE " + tableName + " set name=\\?, updated_at=\\? WHERE ID = \\?"

// 	prep := mock.ExpectPrepare(query)
// 	prep.ExpectExec().WithArgs(dset.Name, dset.UpdatedAt, dset.ID).WillReturnResult(sqlmock.NewResult(12, 1))

// 	a := datasetRepo.NewMysqlDataSetRepository(db)

// 	err = a.Update(context.TODO(), dset)
// 	assert.NoError(t, err)
// }
