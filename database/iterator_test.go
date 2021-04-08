package database

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

type TestingObject struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name" gorm:"size:255;default:'';not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func TestIterator(t *testing.T) {
	tx, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	assert.NoError(t, err)
	err = tx.AutoMigrate(TestingObject{})
	assert.NoError(t, err)

	t.Run("Global", func(t *testing.T) {
		// Not Parallel
		err = tx.Delete(&TestingObject{}, "id > 0").Error
		assert.NoError(t, err)
		objects := []TestingObject{
			{
				Name: "test_iterator_1",
			},
			{
				Name: "test_iterator_2",
			},
			{
				Name: "test_iterator_3",
			},
		}
		err = tx.Create(&objects).Error
		assert.NoError(t, err)

		ts := make([]TestingObject, 2)
		it, err := NewIterator(tx, &ts, 2)
		assert.NoError(t, err)
		for i := 0; true; i++ {
			ok, err := it.Next()
			assert.NoError(t, err)
			if !ok {
				break
			}
			assert.Equal(t, objects[i].ID, ts[it.index].ID)
			assert.Equal(t, objects[i].Name, ts[it.index].Name)
		}
	})
	t.Run("Selected", func(t *testing.T) {
		// Not Parallel
		objects := []TestingObject{
			{
				Name: "test_iterator_search_1",
			},
			{
				Name: "test_iterator_search_2",
			},
			{
				Name: "test_iterator_search_3",
			},
		}
		err = tx.Create(&objects).Error
		assert.NoError(t, err)

		ts := make([]TestingObject, 2)
		it, err := NewIterator(tx.Where("name like ?", "%search%"), &ts, 2)
		assert.NoError(t, err)
		for i := 0; true; i++ {
			ok, err := it.Next()
			assert.NoError(t, err)
			if !ok {
				break
			}
			assert.Equal(t, objects[i].ID, ts[it.index].ID)
			assert.Equal(t, objects[i].Name, ts[it.index].Name)
		}
	})
	t.Run("Empty", func(t *testing.T) {
		// Not Parallel
		err = tx.Delete(&TestingObject{}, "id > 0").Error
		assert.NoError(t, err)

		ts := make([]TestingObject, 2)
		it, err := NewIterator(tx.Where("name like ?", "%search%"), &ts, 2)
		assert.NoError(t, err)
		ok, err := it.Next()
		assert.NoError(t, err)
		assert.False(t, ok)
	})

}
