package main

import (
	"backend/models"
	"backend/services"
	"backend/services/parser"
	"backend/services/search"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Post("/", func(c fiber.Ctx) error {
		if !c.HasBody() {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		req := new(models.Request)

		if err := c.Bind().JSON(req); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if req.URL == "" || req.Type == "" || req.Amount < 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if req.Type != "DFS" && req.Type != "BFS" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if req.Type == "DFS" {
			log.Println("DFS request received: URL=%s, Amount=%d, Selector=%s\n", req.URL, req.Amount, req.Selector)

			rawHTML, err := services.FetchHTMLPage(req.URL)

			if err != nil {
				log.Println("Error fetching HTML page: %v\n", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			DOMTree, err := parser.ParseHTML(rawHTML)

			if err != nil {
				log.Println("Error parsing HTML: %v\n", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			if req.Selector == "" {
				log.Println("No selector provided for DFS search\n")
				return c.SendStatus(fiber.StatusBadRequest)
			}

			sel, err := parser.ParseCSSSelector(req.Selector)

			if err != nil {
				log.Println("Error parsing CSS selector: %v\n", err)
				return c.SendStatus(fiber.StatusBadRequest)
			}

			res, searchlog := search.SearchElementDFS(DOMTree, &sel, req.Amount)

			log.Printf("DFS search result: %v\n", res)
			log.Printf("DFS search log: %v\n", searchlog)

			return c.JSON(map[string]interface{}{
				"result": res.Serialize(),
				"log":    searchlog.Serialize(),
			})

		} else {
			log.Printf("BFS request received: URL=%s, Amount=%d\n", req.URL, req.Amount)
			return c.SendStatus(fiber.StatusOK)

		}

	})

	log.Fatal(app.Listen(":6767"))
}
