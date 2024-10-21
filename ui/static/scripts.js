// scripts.js

// Initialize dropdown with dummy book collections
window.onload = () => {
  fetch("/api/v1/books")
    .then((response) => response.json())
    .then((data) => {
      const dropdown = document.getElementById("bookCollection");
      data.forEach((collection) => {
        const option = document.createElement("option");
        option.value = collection.isbn;
        option.textContent = collection.title;
        dropdown.appendChild(option);
      });
    });
};

document.getElementById("loadShelf").addEventListener("click", () => {
  const selectedCollection = document.getElementById("bookCollection").value;
  const width = parseFloat(document.getElementById("maxWidth").value);

  if (!selectedCollection || !maxWidth) {
    alert("Please select a book collection and enter a valid max width.");
    return;
  }

  fetch(`/api/v1/collection/${selectedCollection}?width=${width}`)
    .then((response) => response.json())
    .then((shelfData) => {
      renderBookshelves(shelfData, width);
    });
});

function renderBookshelves(shelves, maxWidth) {
  const shelvesContainer = document.getElementById("shelvesContainer");
  shelvesContainer.innerHTML = "";

  let shelfCounter = 0;
  shelves.forEach((shelf, index) => {
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

