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
					Username string `gorm:"unique_index;size:30" json:"username" validate:"required,max=30,min=5"`
					Nickname string `gorm:"index:nickname;size:30" json:"nickname"`
					Email    string `gorm:"unique_index;size:320" json:"email"`
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
					Token     string `gorm:"unique_index;size:32" json:"token"`
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
					ID        uint    `gorm:"primary_key"`
					Key       string  `gorm:"size:255"`
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
					Name   string  `json:"name" gorm:"size:255"`
					Target *string `json:"target" gorm:"size:255"`
				}

				type Permission struct {
					ID     uint   `gorm:"primary_key" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
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
		// add UpdateAt column
		{
			ID: "add_updated_at_column_to_tokens",
			Migrate: func(tx *gorm.DB) error {
				type Token struct {
					UpdatedAt time.Time
				}
				return tx.AutoMigrate(&Token{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("tokens").DropColumn("updated_at").Error
			},
		},
		// add RememberMe column
		{
			ID: "add_remember_me_column_to_tokens",
			Migrate: func(tx *gorm.DB) error {
				type Token struct {
					RememberMe bool
				}
				return tx.AutoMigrate(&Token{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Table("tokens").DropColumn("remember_me").Error
			},
		},
		// add images table
		{
			ID: "create_images_table",
			Migrate: func(tx *gorm.DB) error {
				type Image struct {
					ID        uint      `gorm:"primary_key" json:"id"`
					Filename  string    `gorm:"filename,size:2048,unique_index"`
					FilePath  string    `gorm:"filepath,size:2048"`
					UserID    uint      `gorm:"index"`
					CreatedAt time.Time `json:"created_at"`
					UpdatedAt time.Time `json:"updated_at"`
				}

				return tx.AutoMigrate(&Image{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("images").Error
			},
		},
		// add default admin role
		{
			ID: "add_default_admin_role",
			Migrate: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primary_key" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primary_key" json:"id"`
					Name        string  `json:"name" gorm:"size:255"`
					Target      *string `json:"target" gorm:"size:255"`
					Permissions []Permission
				}

				admin := Role{
					Name:   "admin",
					Target: nil,
				}
				err = tx.Create(&admin).Error
				if err != nil {
					return
				}
				perm := Permission{
					RoleID: admin.ID,
					Name:   "all",
				}
				err = tx.Model(&admin).Association("Permissions").Append(perm).Error
				if err != nil {
					return
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primary_key" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primary_key" json:"id"`
					Name        string  `json:"name" gorm:"size:255"`
					Target      *string `json:"target" gorm:"size:255"`
					Permissions []Permission
				}

				var admin Role
				err = tx.Where("name = ? ", "admin").First(&admin).Error // TODO: and target == nil ?
				if err != nil {
					return
				}
				err = tx.Delete(Permission{}, "role_id = ?", admin.ID).Error
				if err != nil {
					return
				}
				err = tx.Delete(&admin).Error
				if err != nil {
					return
				}
				return
			},
		},
		// add problems
		{
			ID: "add_problems",
			Migrate: func(tx *gorm.DB) (err error) {

				type TestCase struct {
					ID uint `gorm:"primary_key" json:"id"`

					ProblemID uint `sql:"index" json:"problem_id" gorm:"not null"`
					Score     uint `json:"score" gorm:"default:0;not null"` // 0 for 平均分配

					InputFileName  string `json:"input_file_name" gorm:"size:255;default:'';not null"`
					OutputFileName string `json:"output_file_name" gorm:"size:255;default:'';not null"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `json:"deleted_at"`
				}

				type Problem struct {
					ID                 uint   `gorm:"primary_key" json:"id"`
					Name               string `sql:"index" json:"name" gorm:"size:255;default:'';not null"`
					Description        string `json:"description"`
					AttachmentFileName string `json:"attachment_file_name" gorm:"size:255;default:'';not null"`
					Public             bool   `json:"public" gorm:"default:false;not null"`
					Privacy            bool   `json:"privacy" gorm:"default:false;not null"`

					MemoryLimit        uint64 `json:"memory_limit" gorm:"default:0;not null;type:bigint"`       // Byte
					TimeLimit          uint   `json:"time_limit" gorm:"default:0;not null"`                     // ms
					LanguageAllowed    string `json:"language_allowed" gorm:"size:255;default:'';not null"`     // E.g.    cpp,c,java,python
					CompileEnvironment string `json:"compile_environment" gorm:"size:2047;default:'';not null"` // E.g.  O2=false
					CompareScriptID    uint   `json:"compare_script_id" gorm:"default:0;not null"`

					TestCases []TestCase `json:"test_cases"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `json:"deleted_at"`
				}
				err = tx.AutoMigrate(&TestCase{}, &Problem{}).Error
				if err != nil {
					return
				}
				if reflect.TypeOf(tx.DB().Driver()).String() != "*sqlite3.SQLiteDriver" {
					err = tx.Model(&TestCase{}).AddForeignKey("problem_id", "problems(id)", "CASCADE", "CASCADE").Error
					if err != nil {
						return
					}
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {
				err = tx.DropTable("test_cases").Error
				if err != nil {
					return
				}
				err = tx.DropTable("problems").Error
				if err != nil {
					return
				}
				return
			},
		},
		{
			// add default problem role
			ID: "add_default_problem_role",
			Migrate: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primary_key" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primary_key" json:"id"`
					Name        string  `json:"name" gorm:"size:255"`
					Target      *string `json:"target" gorm:"size:255"`
					Permissions []Permission
				}

				problemString := "problem"
				problemCreator := Role{
					Name:   "problem_creator",
					Target: &problemString,
				}
				err = tx.Create(&problemCreator).Error
				if err != nil {
					return
				}
				problemPerm := Permission{
					RoleID: problemCreator.ID,
					Name:   "all",
				}
				err = tx.Model(&problemCreator).Association("Permissions").Append(problemPerm).Error
				if err != nil {
					return
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primary_key" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primary_key" json:"id"`
					Name        string  `json:"name" gorm:"size:255"`
					Target      *string `json:"target" gorm:"size:255"`
					Permissions []Permission
				}

				var problemCreator Role
				err = tx.Where("name = ? and target = ? ", "problem_creator", "problem").First(&problemCreator).Error
				if err != nil {
					return
				}
				err = tx.Delete(Permission{}, "role_id = ?", problemCreator.ID).Error
				if err != nil {
					return
				}
				err = tx.Delete(&problemCreator).Error
				if err != nil {
					return
				}
				return
			},
		},
	})
}

func Migrate() {
	m := GetMigration()
	if err := m.Migrate(); err != nil {
		fmt.Printf("Could not migrate: %v", err)
		panic(err)
	}
}
