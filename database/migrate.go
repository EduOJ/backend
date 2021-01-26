package database

import (
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/leoleoasd/EduOJBackend/base"
	"gorm.io/gorm"
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
					ID        uint `gorm:"primaryKey"`
					Level     *int
					Message   string
					Caller    string
					CreatedAt time.Time
				}
				return tx.AutoMigrate(&Log{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("logs")
			},
		},
		// create users table
		{
			ID: "create_users_table",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					ID       uint   `gorm:"primaryKey" json:"id"`
					Username string `gorm:"uniqueIndex;size:30" json:"username" validate:"required,max=30,min=5"`
					Nickname string `gorm:"index;size:30" json:"nickname"`
					Email    string `gorm:"uniqueIndex;size:320" json:"email"`
					Password string `json:"-"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
				}
				return tx.AutoMigrate(&User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		// add tokens table
		{
			ID: "create_tokens_table",
			Migrate: func(tx *gorm.DB) error {
				type User struct {
					ID       uint   `gorm:"primaryKey" json:"id"`
					Username string `gorm:"uniqueIndex;size:30" json:"username" validate:"required,max=30,min=5"`
					Nickname string `gorm:"index;size:30" json:"nickname"`
					Email    string `gorm:"uniqueIndex;size:320" json:"email"`
					Password string `json:"-"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
				}
				type Token struct {
					ID     uint   `gorm:"primaryKey" json:"id"`
					Token  string `gorm:"unique_index;size:32" json:"token"`
					UserID uint
					User
					CreatedAt time.Time `json:"created_at"`
				}
				return tx.AutoMigrate(&Token{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("tokens")
			},
		},
		// add configs table
		{
			ID: "create_configs_table",
			Migrate: func(tx *gorm.DB) error {
				type Config struct {
					ID        uint    `gorm:"primaryKey"`
					Key       string  `gorm:"size:255"`
					Value     *string `gorm:"default:''"` // 可能是空字符串, 因此得是指针
					CreatedAt time.Time
					UpdatedAt time.Time
				}
				return tx.AutoMigrate(&Config{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("configs")
			},
		},
		// add permissions
		{
			ID: "add_permissions",
			Migrate: func(tx *gorm.DB) (err error) {
				type User struct {
					ID       uint   `gorm:"primaryKey" json:"id"`
					Username string `gorm:"uniqueIndex;size:30" json:"username" validate:"required,max=30,min=5"`
					Nickname string `gorm:"index;size:30" json:"nickname"`
					Email    string `gorm:"uniqueIndex;size:320" json:"email"`
					Password string `json:"-"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `gorm:"index" json:"deleted_at"`
				}
				type UserHasRole struct {
					ID     uint `gorm:"primaryKey" json:"id"`
					UserID uint `json:"user_id"`
					User
					RoleID   uint `json:"role_id"`
					TargetID uint `json:"target_id"`
				}

				type Role struct {
					ID     uint    `gorm:"primaryKey" json:"id"`
					Name   string  `json:"name" gorm:"size:255"`
					Target *string `json:"target" gorm:"size:255"`
				}

				type Permission struct {
					ID     uint   `gorm:"primaryKey" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}
				err = tx.AutoMigrate(&UserHasRole{}, &Role{}, &Permission{})
				if err != nil {
					return
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {
				err = tx.Migrator().DropTable("user_has_roles")
				if err != nil {
					return
				}
				err = tx.Migrator().DropTable("permissions")
				if err != nil {
					return
				}
				err = tx.Migrator().DropTable("roles")
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
				return tx.AutoMigrate(&Token{})
			},
			Rollback: func(tx *gorm.DB) error {
				type Token struct {
					RememberMe bool
				}
				return tx.Migrator().DropColumn(&Token{}, "updated_at")
			},
		},
		// add RememberMe column
		{
			ID: "add_remember_me_column_to_tokens",
			Migrate: func(tx *gorm.DB) error {
				type Token struct {
					RememberMe bool
				}
				return tx.AutoMigrate(&Token{})
			},
			Rollback: func(tx *gorm.DB) error {
				type Token struct {
					RememberMe bool
				}
				return tx.Migrator().DropColumn(&Token{}, "remember_me")
			},
		},
		// add images table
		{
			ID: "create_images_table",
			Migrate: func(tx *gorm.DB) error {
				type Image struct {
					ID        uint      `gorm:"primaryKey" json:"id"`
					Filename  string    `gorm:"filename,size:2048,uniqueIndex"`
					FilePath  string    `gorm:"filepath,size:2048"`
					UserID    uint      `gorm:"index"`
					CreatedAt time.Time `json:"created_at"`
					UpdatedAt time.Time `json:"updated_at"`
				}

				return tx.AutoMigrate(&Image{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("images")
			},
		},
		// add default admin role
		{
			ID: "add_default_admin_role",
			Migrate: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primaryKey" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primaryKey" json:"id"`
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
				err = tx.Model(&admin).Association("Permissions").Append(&perm)
				if err != nil {
					return
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primaryKey" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primaryKey" json:"id"`
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
					ID uint `gorm:"primaryKey" json:"id"`

					ProblemID uint `sql:"index" json:"problem_id" gorm:"not null"`
					Score     uint `json:"score" gorm:"default:0;not null"` // 0 for 平均分配

					InputFileName  string `json:"input_file_name" gorm:"size:255;default:'';not null"`
					OutputFileName string `json:"output_file_name" gorm:"size:255;default:'';not null"`

					CreatedAt time.Time  `json:"created_at"`
					UpdatedAt time.Time  `json:"-"`
					DeletedAt *time.Time `json:"deleted_at"`
				}

				type Problem struct {
					ID                 uint   `gorm:"primaryKey" json:"id"`
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
				err = tx.AutoMigrate(&Problem{}, &TestCase{})
				if err != nil {
					return
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {
				err = tx.Migrator().DropTable("test_cases")
				if err != nil {
					return
				}
				err = tx.Migrator().DropTable("problems")
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
					ID     uint   `gorm:"primaryKey" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primaryKey" json:"id"`
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
				err = tx.Model(&problemCreator).Association("Permissions").Append(&problemPerm)
				if err != nil {
					return
				}
				return
			},
			Rollback: func(tx *gorm.DB) (err error) {

				type Permission struct {
					ID     uint   `gorm:"primaryKey" json:"id"`
					RoleID uint   `json:"role_id"`
					Name   string `json:"name" gorm:"size:255"`
				}

				type Role struct {
					ID          uint    `gorm:"primaryKey" json:"id"`
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
		{
			// add sample column
			ID: "add_sample_column",
			Migrate: func(tx *gorm.DB) (err error) {
				type TestCase struct {
					Sample bool `json:"sample" gorm:"default:false;not null"`
				}
				return tx.AutoMigrate(&TestCase{})
			},
			Rollback: func(tx *gorm.DB) (err error) {
				type TestCase struct{}
				return tx.Migrator().DropColumn(&TestCase{}, "sample")
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
