package routes

import (
	"bellesquare-be/internal/handlers"
	repositories "bellesquare-be/internal/repository"
	"bellesquare-be/internal/ws"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRoutes(e *echo.Echo, pool *pgxpool.Pool, secret string, hub *ws.Hub) {
	// 1. Initialize all Repositories
	uRepo := repositories.NewUserRepo(pool)
	tRepo := repositories.NewTeamRepo(pool)
	pRepo := repositories.NewPayoutRepo(pool)
	bRepo := repositories.NewBoardRepo(pool)

	// 2. Setup CORS (for frontend communication) and other defaults
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:4173"}, // SvelteKit dev and preview ports
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// 3. Initialize all Handlers with their respective Repositories
	// public routes
	testRoutes(e)
	authRoutes(e, uRepo, secret)

	// protected routes (require authentication)
	config := echojwt.Config{
		SigningKey:  []byte(secret), // have to cast the secret string to []byte here
		TokenLookup: "header:Authorization:Bearer ,cookie:at",
		Skipper: func(c *echo.Context) bool {
			// Skip auth if the URL ends in /ws
			return strings.HasSuffix(c.Path(), "/ws")
		}, // 🔥 important
	}

	protectedAuthRoutes(e, uRepo, secret, config) // For logout route that needs auth

	apiRoutes := e.Group("/api/v1")
	apiRoutes.Use(echojwt.WithConfig(config)) // Apply JWT middleware to all /api/v1 routes

	userRoutes(apiRoutes, uRepo)
	teamRoutes(apiRoutes, tRepo)
	boardRoutes(apiRoutes, bRepo, pRepo, hub)
	// squareRoutes(apiRoutes, sRepo)
}

/********** THE ROUTES **********/
func testRoutes(e *echo.Echo) {
	e.GET("/hello", handlers.Hello)
	e.GET("/ping", handlers.Ping)
}

func authRoutes(e *echo.Echo, uRepo repositories.UserRepository, secret string) {
	ar := e.Group("/auth")
	ah := handlers.NewAuthHandler(uRepo, secret)
	ar.POST("/login", ah.Login)
	ar.POST("/logout", ah.Logout)
	ar.POST("/register", ah.Register)
}

func protectedAuthRoutes(e *echo.Echo, uRepo repositories.UserRepository, secret string, config echojwt.Config) {
	ar := e.Group("/auth")
	ar.Use(echojwt.WithConfig(config)) // Apply JWT middleware to logout route
	ah := handlers.NewAuthHandler(uRepo, secret)
	ar.GET("/me", ah.Me)
}

func boardRoutes(group *echo.Group, bRepo repositories.BoardRepository, poRepo repositories.PayoutRepository, hub *ws.Hub) {
	br := group.Group("/boards")
	// br.Use(middleware.BodyDump(func(c *echo.Context, reqBody []byte, resBody []byte, err error) {
	// 	fmt.Printf("Request Body: %s\n", string(reqBody))
	// }))
	bh := handlers.NewBoardHandler(bRepo, poRepo, hub)

	// GET /api/v1/boards?user_id=...
	br.GET("", bh.GetUserBoardsHandler)
	br.GET("/:id", bh.GetBoardHandler)
	br.GET("/:id/ws", bh.HandleWS)
	br.GET("/joined", bh.GetJoinedBoards)

	br.Match([]string{http.MethodPost, http.MethodPatch}, "/:id/status", bh.UpdateStatusHandler)

	br.POST("/:id/shuffle", bh.GenerateAxisNumbersHandler)
	br.POST("/new", bh.CreateBoardHandler)
	br.POST("/:id/squares/claim", bh.ClaimSquareHandler)
	br.POST("/:id/payouts", bh.SavePayoutsHandler)
	br.POST("/:id/scores", bh.UpdateScoreHandler)
	br.POST("/:boardId/payment-sent", bh.MarkPaymentSentHandler)
	br.POST("/:boardId/players/:userId/toggle-paid", bh.TogglePaidHandler)
	br.PUT("/:boardId/players/:userId/verify", bh.VerifyPaymentHandler)
}

func userRoutes(group *echo.Group, uRepo repositories.UserRepository) {
	ur := group.Group("/users")
	uh := handlers.NewUserHandler(uRepo)

	ur.POST("/new", uh.AddUserHandler)
	ur.POST("/changepass", uh.UpdatePasswordHandler)
	ur.POST("/upgrade", uh.UpgradeAccount)
	ur.PUT("/profile", uh.UpdateProfile)
}

func teamRoutes(group *echo.Group, tRepo repositories.TeamRepository) {
	tr := group.Group("/teams")
	th := handlers.NewTeamHandler(tRepo)

	tr.GET("", th.GetTeamsHandler)
	tr.GET("/:id", th.GetTeam)
	tr.GET("/counts", th.GetCounts)
	tr.POST("/new", th.AddTeamHandler)
	tr.PUT("/:id", th.UpdateTeamHandler)
	tr.DELETE("/:id", th.DeleteTeamHandler)
}
