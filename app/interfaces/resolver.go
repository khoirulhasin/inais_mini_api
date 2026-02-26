package interfaces

import (
	"github.com/khoirulhasin/untirta_api/app/domains/cams"
	"github.com/khoirulhasin/untirta_api/app/domains/devices"
	"github.com/khoirulhasin/untirta_api/app/domains/drivers"
	"github.com/khoirulhasin/untirta_api/app/domains/drives"
	"github.com/khoirulhasin/untirta_api/app/domains/marker_types"
	"github.com/khoirulhasin/untirta_api/app/domains/markers"
	"github.com/khoirulhasin/untirta_api/app/domains/menus"
	"github.com/khoirulhasin/untirta_api/app/domains/menus2roles"
	"github.com/khoirulhasin/untirta_api/app/domains/profiles"
	"github.com/khoirulhasin/untirta_api/app/domains/roles"
	"github.com/khoirulhasin/untirta_api/app/domains/ships"
	"github.com/khoirulhasin/untirta_api/app/domains/users"
	"github.com/khoirulhasin/untirta_api/app/domains/users2roles"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProfileRepository    profiles.ProfileRepository
	RoleRepository       roles.RoleRepository
	UserRepository       users.UserRepository
	Users2roleRepository users2roles.Users2roleRepository
	MenuRepository       menus.MenuRepository
	Menus2roleRepository menus2roles.Menus2roleRepository
	DeviceRepository     devices.DeviceRepository
	MarkerRepository     markers.MarkerRepository
	ShipRepository       ships.ShipRepository
	DriverRepository     drivers.DriverRepository
	DriveRepository      drives.DriveRepository
	CamRepository        cams.CamRepository
	MarkerTypeRepository marker_types.MarkerTypeRepository
	ShipMongodistory     ships.ShipMongodistory
	ShipMongotory        ships.ShipMongotory
}
