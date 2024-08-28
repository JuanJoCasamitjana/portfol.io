package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/JuanJoCasamitjana/portfol.io/internal/database/sqlite"
	"github.com/JuanJoCasamitjana/portfol.io/internal/model"
	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBname = "dev.db"
var DB *gorm.DB
var Connector *libsql.Connector

const ReplicasDirStr = "./replicas"

var Replicas string

func Remigrate() {
	DB.AutoMigrate(&model.User{}, &model.Article{}, &model.Project{}, &model.Image{}, &model.Gallery{},
		&model.Post{}, &model.Section{}, &model.FollowList{}, &model.Report{}, &model.Tag{})
}

func init() {
	rep, err := os.MkdirTemp(ReplicasDirStr, "libsql-*")
	if err != nil {
		log.Println("replicas directory could not be created")
	}
	Replicas = rep
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	newDBname := os.Getenv("DB_NAME")
	if newDBname != "" {
		DBname = newDBname
	}
	dbPath := filepath.Join(rep, DBname)
	tursoDBUrl := os.Getenv("TURSO_DB_URL")
	tursoDBToken := os.Getenv("TURSO_DB_TOKEN")
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, tursoDBUrl, libsql.WithAuthToken(tursoDBToken), libsql.WithSyncInterval(60*time.Second))
	replicaConnectorCreated := err == nil
	if !replicaConnectorCreated {
		log.Println(err)
	}
	log.Println("connector: ", connector)
	log.Println("is connector created: ", replicaConnectorCreated)
	db := sql.OpenDB(connector)
	Connector = connector
	tursoDSN := fmt.Sprintf("%s?authToken=%s", tursoDBUrl, tursoDBToken)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 300, // Slow SQL threshold
			LogLevel:                  logger.Error,           // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,                  // Include params in the SQL log
			Colorful:                  false,                  // Disable color
		},
	)
	tursoDialector := sqlite.New(sqlite.Config{DriverName: "libsql", DSN: tursoDSN})
	tursoReplicaDialector := sqlite.New(sqlite.Config{DriverName: "libsql", Conn: db})
	dialectorFinal := sqlite.Open(DBname)
	if tursoDBUrl != "" {
		dialectorFinal = tursoDialector
	}
	if replicaConnectorCreated {
		dialectorFinal = tursoReplicaDialector
	}
	DB, err = gorm.Open(dialectorFinal, &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: false, //This ensures data consistency by wrapping atomic operations in transactions
	})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(&model.User{}, &model.Article{}, &model.Project{}, &model.Image{}, &model.Gallery{},
		&model.Post{}, &model.Section{}, &model.FollowList{}, &model.Report{}, &model.Tag{}, &model.Vote{})
	if err != nil {
		log.Fatal(err)
	}
	godotenv.Load()
	ADMIN_USERNAME := os.Getenv("ADMIN_USERNAME")
	ADMIN_PASSWORD := os.Getenv("ADMIN_PASSWORD")
	ADMIN_FULLNAME := os.Getenv("ADMIN_FULLNAME")
	if ADMIN_USERNAME == "" || ADMIN_PASSWORD == "" {
		log.Println("Admin username or password not set, you may use the aplication without admin privileges")
		return
	}
	admin := model.User{
		Username:  ADMIN_USERNAME,
		Authority: model.AUTH_ADMIN,
		FullName:  ADMIN_FULLNAME,
		Profile:   model.Profile{PfPUrl: "/static/default-avatar.png"},
	}
	err = admin.Password.SetPasswordAsHash(ADMIN_PASSWORD)
	if err != nil {
		log.Fatal(err)
	}
	admin_db := model.User{}
	err = DB.Model(admin_db).Where("username = ?", ADMIN_USERNAME).First(&admin_db).Error
	if err == nil {
		log.Println(admin_db.Username)
		log.Println("Admin user already exists")
		return
	}
	if err != gorm.ErrRecordNotFound {
		log.Fatal(err)
	}
	err = DB.Transaction(func(tx *gorm.DB) error {
		tx.Create(&admin)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
