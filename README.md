slackin
=======

Rewrite of [slackin](https://github.com/rauchg/slackin) in go.
Supports multiple organizations.

## Setup

```
make
```

## Configuration

`slackin` takes a json file with your organization's configuration. Check
`config/settings.json.example` for reference.

It assumes that there is a `/banners/<name>.png` file available.

## Usage

Slackin assumes there is a `config/settings.json` file. Once that file is
correctly setup, all you have to do is run `./slackin`.
