const books = [];
const form = document.getElementById("bookForm");
const container = document.getElementById("bookListContainer");

function sendBook(bookData) {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:5000/books", true);

    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.send(bookData);
}

form.addEventListener("submit", function(e) {
    e.preventDefault();

    const book = {
        title: document.getElementById("title").value,
        author: document.getElementById("author").value,
        isbn: document.getElementById("isbn").value,
        price: parseFloat(document.getElementById("price").value).toFixed(2)
    };

    books.push(book);

    var bookData = JSON.stringify(book);
    sendBook(bookData);

    displayBooks();
    form.reset();
});

function displayBooks() {
    if (books.length === 0) {
        container.innerHTML = '<div class="empty-state">No books added yet. Start building your library!</div>';
        return;
    }

    container.innerHTML = books.map(book => `
        <div class="book-item">
            <div class="book-title">${book.title}</div>
            <div class="book-author">by ${book.author}</div>
            <div class="book-details">
                <span><strong>ISBN:</strong> ${book.isbn}</span>
                <span><strong>Price:</strong> $${book.price}</span>
            </div>
        </div>
    `).join("");
}
