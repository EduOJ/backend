package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"gopkg.in/gormigrate.v1"
	"reflect"
	"time"
)

func GetMigration() *gormigrate.Gormigrate {
	return gormigrate.New(base.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// dummy
		{
			ID: "start",
			Migrate: func(tx *gorm.DB) error {
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
		// create logs table
		{
			ID: "create_logs_table",
			Migrate: func(tx *gorm.DB) error {
				type Log struct {
					ID        uint `gorm:"primary_key"`
					Level     *int
					Message   string
					Caller    string
					CreatedAt time.Time
				}
				return tx.AutoMigrate(&Log{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("logs").Error
			},
		},
		// create users table
		{
			ID: "create_users_table",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					ID       uint   `gorm:"primary_key" json:"id"`
					Username string `gorm:"unique_index" json:"username" validate:"required,max=30,min=5"`
					Nickname string `gorm:"index:nickname" json:"nickname"`
					Email    string `gorm:"unique_index" json:"email"`
					Password string `json:"-"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
				}
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("users").Error
			},
		},
		// add tokens table
		{
			ID: "create_tokens_table",
			Migrate: func(tx *gorm.DB) error {
				type Token struct {
					ID        uint   `gorm:"primary_key" json:"id"`
					Token     string `gorm:"unique_index" json:"token"`
					UserID    uint
					CreatedAt time.Time `json:"created_at"`
				}
				err := tx.AutoMigrate(&Token{}).Error
				if err != nil {
					return err
				}
				if reflect.TypeOf(tx.DB().Driver()).String() != "*sqlite3.SQLiteDriver" {
					err = tx.Model(&Token{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
				}
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("tokens").Error
			},
		},
		// add configs table
		{
			ID: "create_configs_table",
			Migrate: func(tx *gorm.DB) error {
				type Config struct {
					ID        uint `gorm:"primary_key"`
					Key       string
					Value     *string `gorm:"default:''"` // 可能是空字符串, 因此得是指针
					CreatedAt time.Time
					UpdatedAt time.Time
				}
				return tx.AutoMigrate(&Config{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("configs").Error
			},
		},
		// add permissions
		{
			ID: "add_permissions",
			Migrate: func(tx *gorm.DB) (err error) {

				type UserHasRole struct {
					ID       uint `gorm:"primary_key" json:"id"`
					UserID   uint `json:"user_id"`
					RoleID   uint `json:"role_id"`
					TargetID uint `json:"target_id"`
				}

				type Role struct {
					ID     uint    `gorm:"primary_key" json:"id"`
					Name   string  `json:"name"`
					Target *string `json:"target"`
				}

				type Permission struct {
					ID     uint   `gorm:"primary_key" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name"`
				}
				err = tx.AutoMigrate(&UserHasRole{}, &Role{}, &Permission{}).Error
				if err != nil {
					return
				}
				if reflect.TypeOf(tx.DB().Driver()).String() != "*sqlite3.SQLiteDriver" {
					err = tx.Model(&UserHasRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE").Error
					if err != nil {
						return
					}
					err = tx.Model(&UserHasRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE").Error
					if err != nil {
						return
					}
					err = tx.Model(&Permission{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE").Error
					if err != nil {
						return
					}
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {
				err = tx.DropTable("user_has_roles").Error
				if err != nil {
					return
				}
				err = tx.DropTable("permissions").Error
				if err != nil {
					return
				}
				err = tx.DropTable("roles").Error
				if err != nil {
					return
				}
				return
			},
		},
		{
			ID: "add_updated_at_column_to_tokens",
			Migrate: func(tx *gorm.DB) error {
				// when table already exists, it just adds fields as columns
				type Token struct {
					UpdatedAt time.Time
				}
				return tx.AutoMigrate(&Token{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("tokens").DropColumn("updated_at").Error
			},
		},
	})
}

func Migrate() {
	/*
		err := base.DB.AutoMigrate(
			&models.Log{}, &models.User{}, &models.Token{}, &models.Config{}, &models.UserHasRole{}, &models.Role{}, &models.Permission{}).Error
		if err != nil {
			fmt.Print(err)
			panic(err)
		}
	*/
	m := GetMigration()
	if err := m.Migrate(); err != nil {
		fmt.Printf("Could not migrate: %v", err)
		panic(err)
	}
}
