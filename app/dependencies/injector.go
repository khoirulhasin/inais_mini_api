package dependencies

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	"github.com/khoirulhasin/untirta_api/app/api/handlers"
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
	"github.com/khoirulhasin/untirta_api/app/generated"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/dbs/mongodb"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/dbs/mongodis"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/dbs/postgres"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/directives"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/middlewares"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/pkg"
	"github.com/khoirulhasin/untirta_api/app/interfaces"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type MyError struct {
	Msg    string
	Status int
}

func (e *MyError) Error() string {
	return e.Msg
}

// Global variable untuk menyimpan handlers agar bisa diakses dari main.go
var GlobalHandlers *Handlers

// Struct untuk menyimpan semua REST handlers
type Handlers struct {
	MarkerHandler *handlers.MarkerHandler
	// Tambahkan handler lain sesuai kebutuhan
}

func Init(r *gin.Engine) http.Handler {
	var connPostgres = postgres.Connect()
	var connMongo = mongodb.Connect()
	var connMongodis = mongodis.Connect()

	// Initialize repositories (sama seperti sebelumnya)
	profileRepository := profiles.NewProfileRepository(connPostgres)
	roleRepository := roles.NewRoleRepository(connPostgres)
	userRepository := users.NewUserRepository(connPostgres)
	users2roleRepository := users2roles.NewUsers2roleRepository(connPostgres)
	menuRepository := menus.NewMenuRepository(connPostgres)
	menus2roleRepository := menus2roles.NewMenus2roleRepository(connPostgres)
	deviceRepository := devices.NewDeviceRepository(connPostgres)
	markerRepository := markers.NewMarkerRepository(connPostgres)
	shipRepository := ships.NewShipRepository(connPostgres)
	driverRepository := drivers.NewDriverRepository(connPostgres)
	driveRepository := drives.NewDriveRepository(connPostgres)
	camRepository := cams.NewCamRepository(connPostgres)
	markerTypeRepository := marker_types.NewMarkerTypeRepository(connPostgres)
	shipMongodistory := ships.NewShipMongodistory(connMongodis)
	shipMongotory := ships.NewShipMongotory(connMongo)

	// Initialize REST API handlers dan simpan ke global variable
	GlobalHandlers = &Handlers{
		MarkerHandler: handlers.NewMarkerHandler(markerRepository),
		// Initialize handler lain
	}

	// GraphQL Configuration (sama seperti sebelumnya)
	c := generated.Config{
		Resolvers: &interfaces.Resolver{
			ProfileRepository:    profileRepository,
			RoleRepository:       roleRepository,
			UserRepository:       userRepository,
			Users2roleRepository: users2roleRepository,
			MenuRepository:       menuRepository,
			Menus2roleRepository: menus2roleRepository,
			DeviceRepository:     deviceRepository,
			MarkerRepository:     markerRepository,
			ShipRepository:       shipRepository,
			DriverRepository:     driverRepository,
			DriveRepository:      driveRepository,
			CamRepository:        camRepository,
			MarkerTypeRepository: markerTypeRepository,
			ShipMongodistory:     shipMongodistory,
			ShipMongotory:        shipMongotory,
		},
	}

	c.Directives.Auth = directives.AuthDirective
	c.Directives.HasRole = directives.HasRoleDirective
	c.Directives.Validate = directives.ValidateDirective

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	h.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.Printf("🔥 Panic caught in resolver: %v", err)
		return gqlerror.Errorf("🔥 Panic caught in resolver: %v", err)
	})

	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	customLog, err := pkg.NewCustomLogger("logs/gorm.log")
	if err != nil {
		panic("failed to create custom logger: " + err.Error())
	}

	middlewares.RegisterCallbacks(connPostgres)
	connPostgres.Logger = customLog
	r.Use(middlewares.CorsMiddleware())
	r.Use(middlewares.GinContextToContextMiddleware())
	r.Use(middlewares.AuthMiddleware(connPostgres))

	return h
}

// Function untuk mendapatkan handlers dari luar package
func GetHandlers() *Handlers {
	return GlobalHandlers
}
