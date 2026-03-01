package server

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/pottrack/backend/internal/config"
	"github.com/pottrack/backend/internal/handler"
	"github.com/pottrack/backend/internal/middleware"
	"github.com/pottrack/backend/internal/migration"
)

// Server wraps Gin engine and dependencies
type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	db     *sqlx.DB
}

func New(engine *gin.Engine, cfg *config.Config) *Server {
	return &Server{engine: engine, cfg: cfg}
}

func (s *Server) Run() error {
	var err error
	s.db, err = sqlx.Connect("postgres", s.cfg.DBUrl)
	if err != nil {
		return err
	}
	s.db.SetConnMaxLifetime(time.Minute * 5)

	// run migrations
	if err := migration.RunMigrations(s.db, "./migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// setup routes
	s.routes()

	log.Printf("starting server on :%s", s.cfg.Port)
	return s.engine.Run(":" + s.cfg.Port)
}

func (s *Server) routes() {
	r := s.engine

	// public
	auth := r.Group("/api/auth")
	hAuth := handler.NewAuthHandler(s.db, s.cfg.JWTSecret)
	auth.POST("/register", hAuth.Register)
	auth.POST("/login", hAuth.Login)

	// protected
	api := r.Group("/api")
	api.Use(middleware.JWT(s.cfg.JWTSecret))

	pot := api.Group("/pot")
	hPot := handler.NewPotHandler(s.db)
	pot.GET("", hPot.List)
	pot.POST("", hPot.Create)
	pot.GET(":id", hPot.Get)
	pot.PATCH(":id", hPot.Update)
	pot.POST(":id/siklus", hPot.StartCycle)

	siklus := api.Group("/siklus")
	hCycle := handler.NewCycleHandler(s.db)
	siklus.POST(":id/arsip", hCycle.Archive)
	siklus.GET(":id/tahap", hCycle.GetStages)
	siklus.PUT(":id/tahap", hCycle.UpdateStages)
	siklus.POST(":id/tugas", hCycle.CreateTask)
	siklus.GET(":id/catatan", hCycle.ListNotes)
	siklus.POST(":id/catatan", hCycle.AddNote)
	siklus.GET(":id/insiden", hCycle.ListIncidents)
	siklus.POST(":id/insiden", hCycle.CreateIncident)
	siklus.GET(":id/panen", hCycle.ListHarvest)
	siklus.POST(":id/panen", hCycle.CreateHarvest)

	tugas := api.Group("/tugas")
	hTask := handler.NewTaskHandler(s.db)
	tugas.GET("", hTask.List)
	tugas.PATCH(":id", hTask.Update)
	tugas.POST(":id/selesai", hTask.Complete)

	insiden := api.Group("/insiden")
	ins := handler.NewIncidentHandler(s.db)
	insiden.PATCH(":id/selesai", ins.Close)
}
