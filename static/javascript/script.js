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
    `${window.location}shorten?short=${short_url_input.value}&long=${long_url_input.value}`
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
        result__short.href = `${data.short_url}`;
        result__short.innerHTML = `1pt.one/${data.short_url}`;
      }
    });
};

// Copy short url to clipboard
result__copyBtn.addEventListener("click", () => {
  result__copyBtn.innerHTML = "Copied";
  navigator.clipboard.writeText(result__short.innerHTML);
});

// Show error if user was redirected because of an invalid short url
window.onload = function () {
  const urlParams = new URLSearchParams(window.location.search);
  const url_not_found = urlParams.get("url_not_found");
  if (url_not_found) {
    window.history.replaceState({}, document.title, "/");
    Swal.fire({
      icon: "error",
      title: `Error`,
      text: `You were redirected because the short url you entered (${url_not_found}) is invalid. Please try again.`,
    });
  }
};
