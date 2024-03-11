package cmd

import (
	"context"
	"fmt"
	"go-img-processing/bootstrap"
	"go-img-processing/internal/controller"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func restCommand(config *bootstrap.Container) *cobra.Command {
	return &cobra.Command{
		Use:   "rest",
		Short: "Run a web server service",
		RunE: func(cmd *cobra.Command, args []string) error {
			ginEngine := gin.Default()
			ginEngine.RedirectTrailingSlash = true
			ginEngine.RemoveExtraSlash = true

			controller.NewHealthController(ginEngine, config)

			port := viper.GetInt("server.port")
			server := &http.Server{
				Addr:    fmt.Sprintf(":%v", port),
				Handler: ginEngine,
			}

			go func() {
				if err := server.ListenAndServe(); err != nil {
					if err == http.ErrServerClosed {
						fmt.Println("server stopped")
					} else {
						fmt.Printf("failed to start server %s", err)
					}
				} else {
					fmt.Printf("server starting at %v", port)
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			<-quit

			fmt.Println("shutdown server...")

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			fmt.Println("shuting down the server...")
			// config.Close

			if err := server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
				fmt.Printf("failed to gracefully shut down the server %s", err)
			}

			return nil
		},
	}
}
