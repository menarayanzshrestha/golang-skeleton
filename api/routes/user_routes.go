package routes

import (
	"github.com/menarayanzshrestha/trello/api/controllers"
	// "github.com/menarayanzshrestha/trello/api/middlewares"
	"github.com/menarayanzshrestha/trello/infrastructure"
)

// UserRoutes struct
type UserRoutes struct {
	logger         infrastructure.Logger
	handler        infrastructure.Router
	userController controllers.UserController
	// authMiddleware middlewares.FirebaseAuthMiddleware
}

func NewUserRoutes(
	logger infrastructure.Logger,
	handler infrastructure.Router,
	userController controllers.UserController,
	// authMiddleware middlewares.FirebaseAuthMiddleware,
) UserRoutes {
	return UserRoutes{
		handler:        handler,
		logger:         logger,
		userController: userController,
		// authMiddleware: authMiddleware,
	}
}

// Setup user routes
func (s UserRoutes) Setup() {
	s.logger.Zap.Info("Setting up user routes")
	// api := s.handler.Group("/api").Use(s.authMiddleware.Handle())
	api := s.handler.Group("/api")
	{
		api.GET("/user", s.userController.GetUser)
		api.GET("/user/:id", s.userController.GetOneUser)
		api.POST("/user", s.userController.SaveUser)
		api.POST("/user/:id", s.userController.UpdateUser)
		api.DELETE("/user/:id", s.userController.DeleteUser)
	}
}
