package postgres

import (
	"fmt"
	"log"

	"github.com/khoirulhasin/untirta_api/app/infrastructures/helpers"
	"github.com/khoirulhasin/untirta_api/app/models"
	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) SeedAll() error {
	log.Println("Starting database seeding...")

	// Untuk PostgreSQL, gunakan CASCADE atau DELETE FROM
	if err := s.clearTables(); err != nil {
		return fmt.Errorf("failed to clear tables: %v", err)
	}

	// Reset auto increment
	if err := s.resetSequences(); err != nil {
		return fmt.Errorf("failed to reset auto increment: %v", err)
	}

	// Seed data in order (respecting foreign keys)
	if err := s.seedRoles(); err != nil {
		return fmt.Errorf("failed to seed roles: %v", err)
	}

	if err := s.seedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	if err := s.seedProfiles(); err != nil {
		return fmt.Errorf("failed to seed profiles: %v", err)
	}

	if err := s.seedUser2Roles(); err != nil {
		return fmt.Errorf("failed to seed users2roles: %v", err)
	}

	if err := s.seedMenus(); err != nil {
		return fmt.Errorf("failed to seed menus: %v", err)
	}

	if err := s.seedMenu2Roles(); err != nil {
		return fmt.Errorf("failed to seed menus2roles: %v", err)
	}

	// Seed new tables
	if err := s.seedShips(); err != nil {
		return fmt.Errorf("failed to seed ships: %v", err)
	}

	if err := s.seedDrivers(); err != nil {
		return fmt.Errorf("failed to seed drivers: %v", err)
	}

	if err := s.seedDevices(); err != nil {
		return fmt.Errorf("failed to seed devices: %v", err)
	}

	if err := s.seedDrives(); err != nil {
		return fmt.Errorf("failed to seed drives: %v", err)
	}

	if err := s.seedMarkerTypes(); err != nil {
		return fmt.Errorf("failed to seed marker types: %v", err)
	}

	// Re-enable foreign key checks
	s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	log.Println("Database seeding completed successfully!")
	return nil
}

// Menggunakan DELETE FROM dengan CASCADE untuk PostgreSQL
func (s *Seeder) clearTables() error {
	// Hapus data dalam urutan yang benar (child tables dulu)
	tables := []string{
		"markers", "marker_types", "drives", "devices", "drivers", "ships",
		"menus2roles", "users2roles", "menus", "profiles", "users", "roles",
	}

	for _, table := range tables {
		// Menggunakan DELETE FROM karena TRUNCATE tidak bisa dengan foreign key
		if err := s.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return err
		}
		log.Printf("Cleared table: %s", table)
	}

	return nil
}

// Reset sequences untuk PostgreSQL
func (s *Seeder) resetSequences() error {
	sequences := []string{
		"roles_id_seq",
		"users_id_seq",
		"users2roles_id_seq",
		"menus_id_seq",
		"menus2roles_id_seq",
		"ships_id_seq",
		"drivers_id_seq",
		"devices_id_seq",
		"drives_id_seq",
		"markers_id_seq",
		"marker_types_id_seq",
	}

	for _, sequence := range sequences {
		// Reset sequence ke 1
		if err := s.db.Exec(fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", sequence)).Error; err != nil {
			// Alternative jika ALTER SEQUENCE gagal
			if err := s.db.Exec(fmt.Sprintf("SELECT setval('%s', 1, false)", sequence)).Error; err != nil {
				log.Printf("Warning: Could not reset sequence %s: %v", sequence, err)
			}
		} else {
			log.Printf("Reset sequence: %s", sequence)
		}
	}

	return nil
}

func (s *Seeder) seedRoles() error {
	roles := []models.Role{
		{
			Code:      "admin",
			Name:      "Administrator",
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
		{
			Code:      "operator",
			Name:      "Operator",
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
		{
			Code:      "owner",
			Name:      "Pemilik",
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
	}

	for _, role := range roles {
		if err := s.db.Create(&role).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d roles", len(roles))
	return nil
}

func (s *Seeder) seedUsers() error {
	password, _ := helpers.HashPassword("Qazwsxedc.123!!")
	users := []models.User{
		{
			Username:  "khoirul",
			Email:     "khoirul.hasin@ppns.ac.id",
			Password:  password,
			LoggedAt:  nil,
			RoleID:    helpers.IntPtr(1),
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
		{
			Username:  "anonym",
			Email:     "anonym@email.com",
			Password:  password,
			LoggedAt:  nil,
			RoleID:    helpers.IntPtr(3),
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
	}

	for _, user := range users {
		if err := s.db.Create(&user).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d users", len(users))
	return nil
}

func (s *Seeder) seedProfiles() error {

	profiles := []models.Profile{
		{
			Name:      "Khoirul Hasin",
			UserID:    1,
			Address:   helpers.StringPtr("Alamat"),
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
		{
			Name:      "Anonymous",
			UserID:    2,
			Address:   helpers.StringPtr("Alamat"),
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
	}

	for _, profile := range profiles {
		if err := s.db.Create(&profile).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d profiles", len(profiles))
	return nil
}

func (s *Seeder) seedUser2Roles() error {
	user2roles := []models.Users2role{
		{
			UserID:    1,
			RoleID:    1,
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
		{
			UserID:    2,
			RoleID:    3,
			CreatedBy: 1, UpdatedBy: helpers.IntPtr(1),
		},
	}

	for _, u2r := range user2roles {
		if err := s.db.Create(&u2r).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d user2roles relationships", len(user2roles))
	return nil
}

func (s *Seeder) seedMenus() error {
	menus := []models.Menu{
		// Level 0 menus (Main categories)
		{Name: "Dashboard", Code: "dashboard", URL: helpers.StringPtr("/portal/dashboard"), Level: helpers.IntPtr(0), MenuID: nil, Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-dashboard"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},

		{Name: "Monitoring", Code: "monitoring", URL: helpers.StringPtr(""), Level: helpers.IntPtr(0), MenuID: nil, Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-desktop"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Transaksi", Code: "transaction", URL: helpers.StringPtr(""), Level: helpers.IntPtr(0), MenuID: nil, Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-money-bill-transfer"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Administrasi", Code: "administration", URL: helpers.StringPtr(""), Level: helpers.IntPtr(0), MenuID: nil, Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-user"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Konfigurasi", Code: "configuration", URL: helpers.StringPtr(""), Level: helpers.IntPtr(0), MenuID: nil, Sequence: helpers.IntPtr(4), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-gear"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},

		// Level 1 menus
		{Name: "Peta", Code: "monitoring-map", URL: helpers.StringPtr("/portal/monitorings/maps"), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(2), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-map"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Perangkat", Code: "transaction-device", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(3), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-laptop-code"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Berlanggan", Code: "transaction-subscription", URL: helpers.StringPtr("/portal/transactions/subscribes"), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(3), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-star"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Pengemudi", Code: "transaction-driver", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(3), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-regular fa-id-card"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Lokasi", Code: "administration-location", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(4), Sequence: helpers.IntPtr(0), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-marker"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kelola Kapal", Code: "administration-ship", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(4), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-car"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kelola Perangkat", Code: "administration-device", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(4), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-mobile"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Langganan", Code: "administration-subscription", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(4), Sequence: helpers.IntPtr(4), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-subscript"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kelola Pengguna", Code: "configuration-user", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(5), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-person"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Atur Menu", Code: "configuration-menu", URL: helpers.StringPtr("/portal/configurations/menus"), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(5), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-burger"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kelola Kamera", Code: "configurations-cams", URL: helpers.StringPtr(""), Level: helpers.IntPtr(1), MenuID: helpers.IntPtr(5), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-camera"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},

		// Level 2 menus
		{Name: "Pasang Perangkat", Code: "transaction-device-trackers2vehicle1", URL: helpers.StringPtr("/portal/transactions/devices/trackers2vehicles"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(7), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-gear"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Fitur Perangkat", Code: "transaction-device-trackers2vehicle", URL: helpers.StringPtr("/portal/transactions/devices/trackers2features"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(7), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-gears"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Mengemudi", Code: "transaction-driver-drive", URL: helpers.StringPtr("transactions/drivers/drives"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(9), Sequence: helpers.IntPtr(0), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-star-of-life"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Negara", Code: "administration-location-country", URL: helpers.StringPtr("/portal/administrations/locations/countries"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(10), Sequence: helpers.IntPtr(0), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-home"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Provinsi", Code: "administration-location-state", URL: helpers.StringPtr("/portal/administrations/locations/states"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(10), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-map"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kota/Kab", Code: "administration-location-city", URL: helpers.StringPtr("/portal/administrations/locations/cities"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(10), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-building"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Perusahaan", Code: "administration-location-company", URL: helpers.StringPtr("/portal/administrations/locations/companies"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(10), Sequence: helpers.IntPtr(4), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-helmet-safety"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Jenis", Code: "administration-ship-type", URL: helpers.StringPtr("/portal/administrations/ships/types"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(11), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-car-side"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Model", Code: "adminstration-ship-model", URL: helpers.StringPtr("/portal/administrations/ships/models"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(11), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-hexagon-nodes"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Merk", Code: "adminstration-ship-make", URL: helpers.StringPtr("/portal/administrations/ships/makes"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(11), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-copyright"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kapal", Code: "adminstration-ship-ship", URL: helpers.StringPtr("/portal/administrations/ships/ships"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(11), Sequence: helpers.IntPtr(4), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-car"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Pengemudi", Code: "administration-ship-driver", URL: helpers.StringPtr("/portal/administrations/ships/drivers"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(11), Sequence: helpers.IntPtr(5), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-id-card"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Pemilik", Code: "administration-ship-owner", URL: helpers.StringPtr("/portal/administrations/ships/owners"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(11), Sequence: helpers.IntPtr(6), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-code-commit"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Seri", Code: "administration-device-seri", URL: helpers.StringPtr("/portal/administrations/devices/series"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(12), Sequence: helpers.IntPtr(0), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-strikethrough"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Status", Code: "administration-device-state", URL: helpers.StringPtr("/portal/administrations/devices/tracker-states"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(12), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-droplet"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Aksi", Code: "administration-device-action", URL: helpers.StringPtr("/portal/administrations/devices/tracker-actions"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(12), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(0), Icon: helpers.StringPtr("fa-fire"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Perangkat", Code: "administration-device-tracker", URL: helpers.StringPtr("/portal/administrations/devices/trackers"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(12), Sequence: helpers.IntPtr(5), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-truck-fast"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Fitur", Code: "administration-subscription-feature", URL: helpers.StringPtr("/portal/administrations/subscriptions/features"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(13), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-feather"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Jenis Harga", Code: "administration-subscription-pricing-type", URL: helpers.StringPtr("/portal/administrations/subscriptions/pricing-types"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(13), Sequence: helpers.IntPtr(2), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-tag"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Jenis Diskon", Code: "administration-subscription-discount-type", URL: helpers.StringPtr("/portal/administrations/subscriptions/discount-types"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(13), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-percent"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Pengguna", Code: "configuration-user-user", URL: helpers.StringPtr("/portal/configurations/users/users"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(14), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-user"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Peran", Code: "configuration-user-role", URL: helpers.StringPtr("/portal/configurations/users/roles"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(14), Sequence: helpers.IntPtr(3), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-face-rolling-eyes"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Kamera", Code: "configurations-cams-cams", URL: helpers.StringPtr("/portal/configurations/cams/cams"), Level: helpers.IntPtr(2), MenuID: helpers.IntPtr(16), Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-camera"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
		{Name: "Peta", Code: "peta", URL: helpers.StringPtr("/map"), Level: helpers.IntPtr(0), MenuID: nil, Sequence: helpers.IntPtr(1), Show: helpers.IntPtr(1), Icon: helpers.StringPtr("fa-map"), CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},
	}

	for _, menu := range menus {
		if err := s.db.Create(&menu).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d menus", len(menus))
	return nil
}

func (s *Seeder) seedMenu2Roles() error {
	menu2roles := []models.Menus2role{
		{RoleID: 1, MenuID: 4, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},  // Administrasi
		{RoleID: 1, MenuID: 5, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},  // Konfigurasi
		{RoleID: 1, MenuID: 10, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)}, // Lokasi
		{RoleID: 1, MenuID: 21, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)}, // Provinsi
		{RoleID: 1, MenuID: 20, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)}, // Negara
		{RoleID: 1, MenuID: 15, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)}, // Atur Menu
		{RoleID: 1, MenuID: 36, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)}, // Pengguna
		{RoleID: 1, MenuID: 3, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},  // Transaksi
		{RoleID: 1, MenuID: 2, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},  // Monitoring
		{RoleID: 1, MenuID: 6, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},  // Peta
		{RoleID: 1, MenuID: 1, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)},  // Dashboard
		{RoleID: 1, MenuID: 40, CreatedBy: 1, UpdatedBy: helpers.IntPtr(1)}, // Peta
	}

	for _, m2r := range menu2roles {
		if err := s.db.Create(&m2r).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d menu2roles relationships", len(menu2roles))
	return nil
}

func (s *Seeder) seedShips() error {
	ships := []models.Ship{
		{
			Name:        "Uncategorized",
			Number:      helpers.StringPtr("00000000000"),
			Description: nil,
			CreatedBy:   1,
			UpdatedBy:   helpers.IntPtr(1),
		},
	}

	for _, ship := range ships {
		if err := s.db.Create(&ship).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d ships", len(ships))
	return nil
}

func (s *Seeder) seedDrivers() error {
	drivers := []models.Driver{
		{
			Name:             "Anonym",
			NumberIdentifier: "000000000000000000",
			Address:          helpers.StringPtr("Alamat"),
			CreatedBy:        1,
			UpdatedBy:        helpers.IntPtr(1),
		},
	}

	for _, driver := range drivers {
		if err := s.db.Create(&driver).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d drivers", len(drivers))
	return nil
}

func (s *Seeder) seedDevices() error {
	devices := []models.Device{
		{
			Imei:      "0000000000000000",
			Name:      "Dummy",
			OwnerID:   helpers.IntPtr(2),
			ShipID:    helpers.IntPtr(1),
			CreatedBy: 1,
			UpdatedBy: helpers.IntPtr(1),
		},
	}

	for _, device := range devices {
		if err := s.db.Create(&device).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d devices", len(devices))
	return nil
}

func (s *Seeder) seedDrives() error {
	drives := []models.Drive{
		{
			DriverID:    1,
			ShipID:      1,
			Description: helpers.StringPtr("Deskripsi"),
			CreatedBy:   1,
			UpdatedBy:   helpers.IntPtr(1),
		},
	}

	for _, drive := range drives {
		if err := s.db.Create(&drive).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d drives", len(drives))
	return nil
}

func (s *Seeder) seedMarkerTypes() error {
	markerTypes := []models.MarkerType{
		{
			Name:      "Uncategorized",
			Icon:      "🞈",
			CreatedBy: 1,
			UpdatedBy: helpers.IntPtr(1),
		},
		{
			Name:      "Sebaran Ikan",
			Icon:      "🐟",
			CreatedBy: 1,
			UpdatedBy: helpers.IntPtr(1),
		},
	}

	for _, markerType := range markerTypes {
		if err := s.db.Create(&markerType).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d types", len(markerTypes))
	return nil
}
