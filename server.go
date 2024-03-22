package main

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/mtslzr/pokeapi-go"
)

func main() {

	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New())

	app.Static("/css", "./css")

	app.Get("/", func(c *fiber.Ctx) error {
		l, _ := pokeapi.Resource("pokemon")

		for i, p := range l.Results {
			l.Results[i].Name = capitalize(p.Name)
		}

		return c.Render("index", fiber.Map{
			"Results": l.Results,
		})
	})

	app.Get("/pokemon/:name", func(c *fiber.Ctx) error {
		name := strings.ToLower(c.Params("name"))
		p, _ := pokeapi.Pokemon(name)

		for i, a := range p.Abilities {
			p.Abilities[i].Ability.Name = capitalize(a.Ability.Name)
		}

		return c.Render("pokemon-detail", fiber.Map{
			"ImageUrl":  p.Sprites.FrontDefault,
			"Name":      capitalize(p.Name),
			"Height":    p.Height,
			"Weight":    p.Weight,
			"Hp":        p.Stats[0].BaseStat,
			"Abilities": p.Abilities,
			"Types":     p.Types,
		})
	})

	app.Get("/search", func(c *fiber.Ctx) error {
		q := strings.ToLower(c.Query("q"))
		p, err := pokeapi.Pokemon(q)

		if err != nil {
			return c.Render("pokemon-detail", fiber.Map{
				"Error": "Oops! There seems to be no Pokemon by that name in our Pokedex.",
			})
		}

		for i, a := range p.Abilities {
			p.Abilities[i].Ability.Name = capitalize(a.Ability.Name)
		}

		return c.Render("pokemon-detail", fiber.Map{
			"ImageUrl":  p.Sprites.FrontDefault,
			"Name":      capitalize(p.Name),
			"Height":    p.Height,
			"Weight":    p.Weight,
			"Hp":        p.Stats[0].BaseStat,
			"Abilities": p.Abilities,
			"Types":     p.Types,
		})
	})

	log.Fatal(app.Listen(":3000"))
}

func capitalize(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
