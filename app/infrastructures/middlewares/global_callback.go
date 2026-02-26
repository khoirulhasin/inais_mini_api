package middlewares

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func RegisterCallbacks(db *gorm.DB) {
	// Callback untuk operasi CREATE

	db.Callback().Create().Before("gorm:create").Register("set_created_by", func(db *gorm.DB) {
		if db.Statement.Model != nil && db.Statement.Context != nil {
			// Asumsikan user ID tersedia dari konteks (misalnya dari middleware)
			user := ForContext(db.Statement.Context) // Fungsi ini harus diimplementasikan

			if user != nil {
				db.Statement.SetColumn("created_by", &user.ID)
				db.Statement.SetColumn("updated_by", &user.ID) // Untuk create, updated_by sama dengan created_by
			}

		}
	})

	// Callback untuk operasi UPDATE
	db.Callback().Update().Before("gorm:update").Register("set_updated_by", func(db *gorm.DB) {
		if db.Statement.Model != nil && db.Statement.Model != "Query" {
			// Asumsikan user ID tersedia dari konteks
			user := ForContext(db.Statement.Context) // Fungsi ini harus diimplementasikan // Fungsi ini harus diimplementasikan
			if user != nil {
				db.Statement.SetColumn("updated_by", &user.ID)
			}
		}
	})

	db.Callback().Delete().Before("gorm:delete").Register("set_deleted_by", func(db *gorm.DB) {
		// Pastikan operasi adalah soft delete (bukan hard delete)
		curTime := db.Statement.DB.NowFunc().UnixMilli()
		user := ForContext(db.Statement.Context) // Fungsi ini harus diimplementasikan // Fungsi ini harus diimplementasikan
		deletedBy := 0

		if user != nil {
			deletedBy = user.ID
		}

		db.Statement.AddClause(clause.Update{})

		db.Statement.AddClause(clause.Set{
			{Column: clause.Column{Name: "deleted_at"}, Value: curTime},
			{Column: clause.Column{Name: "deleted_by"}, Value: deletedBy},
		})

		db.Statement.SetColumn("deleted_at", curTime)
		db.Statement.SetColumn("deleted_by", deletedBy)

		db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: "deleted_at", Value: 0},
		}})

		db.Statement.Build(
			clause.Update{}.Name(),
			clause.Set{}.Name(),
			clause.Where{}.Name(),
		)

		return
	})

}
