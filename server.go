package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"github.com/mtslzr/pokeapi-go"
)

const pageSize = 20

func main() {

	engine := html.New("./templates", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New())

	app.Static("/css", "./css")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/pokemon-list/0")
	})

	app.Get("/pokemon-list/:offset", func(c *fiber.Ctx) error {
		offset, _ := strconv.Atoi(c.Params("offset"))
		l, _ := pokeapi.Resource("pokemon", offset, pageSize)

		for i, p := range l.Results {
			l.Results[i].Name = capitalize(p.Name)
		}

		// Compute the offsets for the previous and next page.
		nextOffset := offset + len(l.Results)
		if nextOffset >= l.Count {
			nextOffset = -1
		}

		prevOffset := offset - pageSize
		if prevOffset < 0 {
			prevOffset = -1
		}

		return c.Render("index", fiber.Map{
			"Results":    l.Results,
			"NextOffset": nextOffset,
			"PrevOffset": prevOffset,
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
