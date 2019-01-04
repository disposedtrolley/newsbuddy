# newsbuddy

`newsbuddy` is utility to help generate the bulk of my Friday newsletters which I distribute at work.

The goal is to

1. generate an [mjml](https://mjml.io/) template given a plain-text file containing metadata, welcome message, and links to articles
2. invoke the `mjml` command-line tool to generate the HTML file
3. copy the contents of the HTML file to the clipboard

Most of the tooling will be written in Go. Orchestration will be performed using a `Makefile`.

## Usage

From a `.toml` source file, `newsbuddy` will generate an `.mjml` file according to a defined template. The structure of the newsletter is fairly simple, comprising of a header containing the issue number and date, and a body with the welcome message and articles. Source data is expressed in a [TOML](https://github.com/toml-lang/toml) file which is parsed and transformed into the output MJML file.

### Source `.toml`

The source file should be in the following format:

```
[metadata]
title = ""
no = <>
date = ""
welcome = """

"""   # Multiline welcome message

[[articles]]
url = ""
type = "<TEXT | VIDEO>"
cat = "<GENERAL | DEV | RECIPE>"

[[articles]]
url = ""
type = "<TEXT | VIDEO>"
cat = "<GENERAL | DEV | RECIPE>"

# Add as many articles as you wish!
```

### `mjml` template

`template.mjml` contains the empty template used for the newsletter. We use Go's [`text/template`](https://dlintw.github.io/gobyexample/public/text-template.html) package to populate the template.

## todo

- [x] Decide on a config file structure - this is used to define the metadata, welcome message, and links for a specific issue
- [x] Implement text parser for the above
- [ ] Add `Makefile` to orchestrate template generation and invoking `mjml`