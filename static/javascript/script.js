const shorten_btn = document.getElementById("shorten_btn");
const long_url_input = document.getElementById("long_url_input");
const short_url_input = document.getElementById("short_url_input");
const result_text = document.getElementById("result_text");

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

shorten_btn.onclick = function () {
  fetch(
    `http://127.0.0.1:5000/shorten?short=${short_url_input.value}&long=${long_url_input.value}`
  )
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      result_text.innerHTML = data.message;
      Swal.fire({
        icon: "error",
        title: "Oops...",
        text: "Something went wrong!",
        footer: '<a href="">Why do I have this issue?</a>',
      });
    });
};
