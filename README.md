# 1pt.one - A minimalistic URL shortener

A minimalistic URL shortener. Create a short URL in seconds and track the number of visits with the website or the api.
You can use the service with the [website 1pt.one](https://1pt.one/) or the [CLI tool](https://github.com/Jeusto/1pt.one-cli).

![demo](static/assets/demo.gif)

## API

### `/shorten`

#### Method: `GET`

| Parameter | Description                                 | Example              |
| --------- | ------------------------------------------- | -------------------- |
| `long`    | **Required** - The long URL to shorten      | `https://github.com` |
| `short`   | **Optional** - The short version of the url | `gth`                |

#### Example Response

`1pt.one/shorten?short=gth&long=https://github.com`

```json
{
  "message": "Successfully added short url!",
  "short_url": "gth",
  "long_url": "https://github.com",
  "file_content": ""
}
```

### `/retrieve`

#### Method: `GET`

| Parameter | Description                                 | Example |
| --------- | ------------------------------------------- | ------- |
| `short`   | **Required** - The short version of the url | `gth`   |

#### Example Response

```json
{
  "short_url": "gth",
  "long_url": "https://github.com",
  "created_at": "18/09/2021 20:26:17",
  "number_of_visits": 2,
  "file_content": ""
}
```

### `/status`

#### Method: `GET`

#### Example Response

```json
{
  "status": 200,
  "message": "Api is live. Read the documentation at https://github.com/Jeusto/1pt.one"
}
```
