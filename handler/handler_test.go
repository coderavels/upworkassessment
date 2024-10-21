package handler

import (
	"strconv"
	"strings"
	"testing"

	"github.com/coderavels/upworkassessment/client"
	"github.com/coderavels/upworkassessment/fixtures"
	mocks "github.com/coderavels/upworkassessment/handler/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type handlerMocks struct {
	assessClient *mocks.MockAssessClient
}

func setupTest(t *testing.T) (Handler, handlerMocks) {
	ctrl := gomock.NewController(t)
	mockAssessClient := mocks.NewMockAssessClient(ctrl)

	mocks := handlerMocks{
		assessClient: mockAssessClient,
	}

	handler := Handler{
		assessClient: mocks.assessClient,
	}

	return handler, mocks
}

func TestHandler_getBookCollection(t *testing.T) {
	t.Run("fetches correct book collection", func(t *testing.T) {
		handler, mocks := setupTest(t)

		mocks.assessClient.EXPECT().GetBook(fixtures.BookCrimeAndPunishment.ISBN).Return(fixtures.BookCrimeAndPunishment, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.Book1984.ISBN).Return(fixtures.Book1984, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookToKillAMockingBird.ISBN).Return(fixtures.BookToKillAMockingBird, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookBraveNewWorld.ISBN).Return(fixtures.BookBraveNewWorld, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookTheGreatGatsby.ISBN).Return(fixtures.BookTheGreatGatsby, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookPrideAndPrejudice.ISBN).Return(fixtures.BookPrideAndPrejudice, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookTheCatcherInTheRye.ISBN).Return(fixtures.BookTheCatcherInTheRye, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookMobyDick.ISBN).Return(fixtures.BookMobyDick, nil)
		mocks.assessClient.EXPECT().GetBook(fixtures.BookThePictureOfDorianGray.ISBN).Return(fixtures.BookThePictureOfDorianGray, nil)

		booksInCollection, err := handler.getBookCollectionFromClient(fixtures.BookCrimeAndPunishment.ISBN)
		assert.NoError(t, err)
		assert.Equal(t, 9, len(booksInCollection))

		for _, b := range booksInCollection {
			assert.Contains(t, []string{
				fixtures.BookCrimeAndPunishment.ISBN,
				fixtures.Book1984.ISBN,
				fixtures.BookToKillAMockingBird.ISBN,
				fixtures.BookBraveNewWorld.ISBN,
				fixtures.BookTheGreatGatsby.ISBN,
				fixtures.BookPrideAndPrejudice.ISBN,
				fixtures.BookTheCatcherInTheRye.ISBN,
				fixtures.BookMobyDick.ISBN,
				fixtures.BookThePictureOfDorianGray.ISBN,
			}, b.ISBN)
		}
	})
}

func TestOrganiseCollectionInShelves(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		organisedCollection, err := organiseCollectionInShelves([]client.BookDetails{
			fixtures.BookCrimeAndPunishment,
			fixtures.Book1984,
			fixtures.BookToKillAMockingBird,
			fixtures.BookBraveNewWorld,
			fixtures.BookTheGreatGatsby,
			fixtures.BookPrideAndPrejudice,
			fixtures.BookTheCatcherInTheRye,
			fixtures.BookMobyDick,
			fixtures.BookThePictureOfDorianGray,
		}, 50)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(organisedCollection))
		assert.Equal(t, 5, len(organisedCollection[0]))
		firstRowSum := 0

		for _, b := range organisedCollection[0] {
			bookWidth, err := strconv.Atoi(strings.TrimSuffix(b.Width, "cm"))
			assert.NoError(t, err)
			firstRowSum += bookWidth
		}
		assert.Equal(t, 50, firstRowSum)

		secondRowSum := 0

		for _, b := range organisedCollection[1] {
			bookWidth, err := strconv.Atoi(strings.TrimSuffix(b.Width, "cm"))
			assert.NoError(t, err)
			secondRowSum += bookWidth
		}
		assert.Equal(t, 33, secondRowSum)
	})

}
