// scripts.js

// Initialize dropdown with dummy book collections
window.onload = () => {
  displayLoading()
  fetch("/api/v1/books")
    .then((response) => response.json())
    .then((data) => {
      hideLoading();
      const dropdown = document.getElementById("bookCollection");
      data.forEach((collection) => {
        const option = document.createElement("option");
        option.value = collection.isbn;
        option.textContent = collection.title;
        dropdown.appendChild(option);
      });
    })
    .catch((error) => {
      hideLoading();
      console.error(error);
      alert("An error occurred while fetching list of books");
    });
};

document.getElementById("loadShelf").addEventListener("click", () => {
  const selectedCollection = document.getElementById("bookCollection").value;
  const width = parseFloat(document.getElementById("maxWidth").value);
  if (!selectedCollection || !width) {
    alert("Please select a book and enter a valid shelf width.");
    return;
  }

  displayLoading();
  fetch(`/api/v1/collection/${selectedCollection}?width=${width}`)
    .then((response) => response.json())
    .then((shelfData) => {
      hideLoading();
      renderBookshelves(shelfData, width);
    })
    .catch((error) => {
      hideLoading();
      console.error(error);
      alert("An error occurred while fetching the book collection.");
    });
});

// selecting loading div
const loader = document.querySelector("#loading");

// showing loading
function displayLoading() {
	loader.classList.add("display");
}

// hide loading
function hideLoading() {
loader.classList.remove("display");
}

function renderBookshelves(shelves, maxWidth) {
  const shelvesContainer = document.getElementById("shelvesContainer");
    shelvesContainer.innerHTML = `Total Shelves: ${shelves.total_shelves}`;

  let shelfCounter = 0;
  shelves.shelves.forEach((shelf, index) => {
    console.log(shelf);
    if (index % 5 === 0) {
      const bookshelf = document.createElement("div");
      bookshelf.className = "bookshelf";
      shelvesContainer.appendChild(bookshelf);
      shelfCounter++;
    }

    const shelfElement = document.createElement("div");
    shelfElement.className = "shelf";

    shelf.forEach((book) => {
      const bookElement = document.createElement("div");
      bookElement.className = "book";
      bookElement.innerHTML = book.title;

      // Adjust book width relative to the total width of books, not the max shelf width
      const relativeWidth = (parseFloat(book.width) / maxWidth) * 100;
      bookElement.style.width = relativeWidth + "%";

      const randomColor = `rgb(${Math.floor(Math.random() * 256)}, ${Math.floor(
        Math.random() * 256
      )}, ${Math.floor(Math.random() * 256)})`;
      bookElement.style.backgroundColor = randomColor;

      const tooltip = document.createElement("div");
      tooltip.className = "tooltip";
      tooltip.textContent = book.title;

      bookElement.appendChild(tooltip);
      shelfElement.appendChild(bookElement);
    });

    document
      .querySelectorAll(".bookshelf")
      [shelfCounter - 1].appendChild(shelfElement);
  });
}

