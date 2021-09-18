// Dom elements
const shorten_btn = document.getElementById("shorten_btn");
const long_url_input = document.getElementById("long_url_input");
const short_url_input = document.getElementById("short_url_input");

const result__inner = document.getElementById("result__inner");
const result__short = document.getElementById("result__short");
const result__long = document.getElementById("result__long");
const result__copyBtn = document.getElementById("result__copyBtn");

// Click submit button if enter key pressed
long_url_input.addEventListener("keyup", function (event) {
  event.preventDefault();
  if (event.keyCode == 13) {
    shorten_btn.click();
  }
});
short_url_input.addEventListener("keyup", function (event) {
  event.preventDefault();
  if (event.keyCode == 13) {
    shorten_btn.click();
  }
});

// Send request to server to shorten url
shorten_btn.onclick = function () {
  fetch(
    `http://127.0.0.1:5000/shorten?short=${short_url_input.value}&long=${long_url_input.value}`
  )
    .then((response) => response.json())
    .then((data) => {
      // Show modal if there's an error
      if (data.status === 400) {
        Swal.fire({
          icon: "error",
          title: `Error ${data.status}`,
          text: `${data.message}`,
        });
        // Display short url otherwise
      } else {
        result__inner.style.display = "flex";
        result__copyBtn.innerHTML = "Copy";
        result__long.innerHTML = `${data.long_url}`;
        result__short.href = `${data.short_url}`;
        result__short.innerHTML = `1pt.one/${data.short_url}`;
      }
    });
};

result__copyBtn.addEventListener("click", () => {
  result__copyBtn.innerHTML = "Copied";
});
