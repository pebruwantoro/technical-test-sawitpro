package main

import (
	"net/http"
	"os"

	"github.com/pebruwantoro/technical-test-sawitpro/generated"
	"github.com/pebruwantoro/technical-test-sawitpro/handler"
	"github.com/pebruwantoro/technical-test-sawitpro/repository"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)

	e.GET("/swagger.json", func(c echo.Context) error {
		spec, err := generated.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error loading swagger spec")
		}
		return c.JSON(http.StatusOK, spec)
	})

	e.GET("/swagger", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
		<!DOCTYPE html>
		<html lang="en">
		  <head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>API Documentation</title>
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.14.0/swagger-ui.css" />
		  </head>
		  <body>
			<div id="swagger-ui"></div>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.14.0/swagger-ui-bundle.js"></script>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.14.0/swagger-ui-standalone-preset.js"></script>
			<script>
			  const ui = SwaggerUIBundle({
				url: "/swagger.json",  // Swagger spec JSON generated by oapi-codegen
				dom_id: '#swagger-ui',
				deepLinking: true,
				presets: [
				  SwaggerUIBundle.presets.apis,
				  SwaggerUIBundle.SwaggerUIStandalonePreset
				],
				layout: "BaseLayout"
			  });
			</script>
		  </body>
		</html>
		`)
	})

	e.Use(middleware.Logger())
	e.GET(generated.GetSwagger())
	e.Logger.Fatal(e.Start(":8080"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
	}

	return handler.NewServer(opts)
}
