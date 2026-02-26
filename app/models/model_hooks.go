package models

// RegisterCallbacks mendaftarkan global callback untuk GORM

//We want our ids to be uuids, so we define that here

// func (mod *Question) BeforeCreate(scope *gorm.Scope) error {
// 	uuid := uuid.NewV4()
// 	return scope.SetColumn("id", uuid.String())
// }

// func (mod *QuestionOption) BeforeCreate(scope *gorm.Scope) error {
// 	uuid := uuid.NewV4()
// 	return scope.SetColumn("id", uuid.String())
// }

// func (mod *Answer) BeforeCreate(scope *gorm.Scope) error {
// 	uuid := uuid.NewV4()
// 	return scope.SetColumn("id", uuid.String())
// }
