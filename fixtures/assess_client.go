package fixtures

import "github.com/coderavels/upworkassessment/client"

// books
var (
	BookCrimeAndPunishment = client.BookDetails{
		Title:   "Crime and Punishment",
		ISBN:    "978-0679734505",
		Width:   "12cm",
		Related: "978-0451524935", // Book1984
	}

	Book1984 = client.BookDetails{
		Title:   "1984",
		ISBN:    "978-0451524935",
		Width:   "8cm",
		Related: "978-0061120084", // BookToKillAMockingBird
	}

	BookToKillAMockingBird = client.BookDetails{
		Title:   "To Kill a Mockingbird",
		ISBN:    "978-0061120084",
		Width:   "9cm",
		Related: "978-0060850524", // BookBraveNewWorld
	}

	BookBraveNewWorld = client.BookDetails{
		Title:   "Brave New World",
		ISBN:    "978-0060850524",
		Width:   "11cm",
		Related: "978-0743273565", // BookTheGreatGatsby
	}

	BookTheGreatGatsby = client.BookDetails{
		Title:   "The Great Gatsby",
		ISBN:    "978-0743273565",
		Width:   "7cm",
		Related: "978-0141439518", // BookPrideAndPrejudice
	}

	BookPrideAndPrejudice = client.BookDetails{
		ISBN:    "978-0141439518",
		Title:   "Pride and Prejudice",
		Width:   "8cm",
		Related: "978-0316769488", // BookTheCatcherInTheRye
	}

	BookTheCatcherInTheRye = client.BookDetails{
		Title:   "The Catcher in the Rye",
		ISBN:    "978-0316769488",
		Width:   "10cm",
		Related: "978-1853260087", // BookMobyDick
	}

	BookMobyDick = client.BookDetails{
		Title:   "Moby-Dick",
		ISBN:    "978-1853260087",
		Width:   "9cm",
		Related: "978-0141439570", // BookThePictureOfDorianGray
	}

	BookThePictureOfDorianGray = client.BookDetails{
		Title:   "The Picture of Dorian Gray",
		ISBN:    "978-0141439570",
		Width:   "9cm",
		Related: "978-0679734505", // BookCrimeAndPunishment
	}
)
