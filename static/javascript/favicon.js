const { to, set, timeline } = gsap;

function validURL(str) {
  let pattern = new RegExp(
    "^(https?:\\/\\/)?" +
      "((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|" +
      "((\\d{1,3}\\.){3}\\d{1,3}))" +
      "(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*" +
      "(\\?[;&a-z\\d%_.~+=-]*)?" +
      "(\\#[-a-z\\d_]*)?$",
    "i"
  );
  return !!pattern.test(str);
}

function delay(fn, ms) {
  let timer = 0;
  return function (...args) {
    clearTimeout(timer);
    timer = setTimeout(fn.bind(this, ...args), ms || 0);
  };
}

document.querySelectorAll(".url-input").forEach((elem) => {
  let icon = elem.querySelector(".icon"),
    favicon = icon.querySelector(".favicon"),
    input = elem.querySelector("input");

  input.addEventListener(
    "input",
    delay((e) => {
      let bool = input.value.length,
        valid = validURL(input.value);
      to(elem, {
        "--clear-scale": bool ? 1 : 0,
        duration: bool ? 0.5 : 0.15,
        ease: bool ? "elastic.out(1, .7)" : "none",
      });
      to(elem, {
        "--clear-opacity": bool ? 1 : 0,
        duration: 0.15,
      });
      to(elem, {
        "--icon-offset": valid ? "24px" : "0px",
        duration: 0.15,
        delay: valid ? 0 : 0.2,
      });
      if (valid) {
        if (favicon.querySelector("img")) {
          favicon.querySelector("img").src =
            "https://f1.allesedv.com/64/" + input.value;
          return;
        }
        let img = new Image();
        img.onload = () => {
          favicon.appendChild(img);
          to(elem, {
            "--favicon-scale": 1,
            duration: 0.5,
            delay: 0.2,
            ease: "elastic.out(1, .7)",
          });
        };
        img.src = "https://f1.allesedv.com/64/" + input.value;
      } else {
        if (favicon.querySelector("img")) {
          to(elem, {
            "--favicon-scale": 0,
            duration: 0.15,
            onComplete() {
              favicon.querySelector("img").remove();
            },
          });
        }
      }
    }, 250)
  );
});

function getPoint(point, i, a, smoothing) {
  let cp = (current, previous, next, reverse) => {
      let p = previous || current,
        n = next || current,
        o = {
          length: Math.sqrt(
            Math.pow(n[0] - p[0], 2) + Math.pow(n[1] - p[1], 2)
          ),
          angle: Math.atan2(n[1] - p[1], n[0] - p[0]),
        },
        angle = o.angle + (reverse ? Math.PI : 0),
        length = o.length * smoothing;
      return [
        current[0] + Math.cos(angle) * length,
        current[1] + Math.sin(angle) * length,
      ];
    },
    cps = cp(a[i - 1], a[i - 2], point, false),
    cpe = cp(point, a[i - 1], a[i + 1], true);
  return `C ${cps[0]},${cps[1]} ${cpe[0]},${cpe[1]} ${point[0]},${point[1]}`;
}

function getPath(x, smoothing) {
  return [
    [2, 2],
    [12 - x, 12 + x],
    [22, 22],
  ].reduce(
    (acc, point, i, a) =>
      i === 0
        ? `M ${point[0]},${point[1]}`
        : `${acc} ${getPoint(point, i, a, smoothing)}`,
    ""
  );
}
