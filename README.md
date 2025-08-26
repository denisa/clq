# clq — Changelog validation and query tool

![GitHub Release Date](https://img.shields.io/github/release-date/denisa/clq?color=blue)
[![version](https://img.shields.io/github/v/release/denisa/clq?include_prereleases&sort=semver)](https://github.com/denisa/clq/releases)
[![Docker Image Version](https://img.shields.io/docker/v/denisa/clq?label=docker%20tag&sort=semver)](https://hub.docker.com/repository/docker/denisa/clq)

[![ci](https://github.com/denisa/clq/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/denisa/clq/actions/workflows/ci.yaml?query=branch%3Amain)
[![coverage status](https://coveralls.io/repos/github/denisa/clq/badge.svg?branch=main)](https://coveralls.io/github/denisa/clq?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/denisa/clq)](https://goreportcard.com/report/github.com/denisa/clq)
[![Super-Linter](https://github.com/denisa/clq/actions/workflows/linter.yaml/badge.svg)](https://github.com/denisa/clq/actions/workflows/linter.yaml)

## usage

clq always validates the complete changelog, stopping at its first error.
If a query is given, clq then queries the changelog and returns the query result.
clq handles standard input — when no arguments are present or an argument is "-" — or any number of files.

clq exits with a status of 0 if all files are valid, and with a non-zero status
if any file fails to validate. It writes to standard output the result of the query if a query was given.

clq writes validation error to standard error.

When processing multiple files, clq prefixes every line on standard out and standard error with the filename.

```text
Usage: clq { options } path_to_changelog.md

Options are:
  -changeMap name
      name of a file defining the mapping from change kind to semantic version change
  -output format
      the format to apply to the result of a (complex) query. Supports `json` and `md` (markdown); defaults to `json`
  -query string
      A query to extract information out of the change log
  -release
      Enable release-mode validation
  -with-filename
      Always print filename headers with output lines
```

Example:

- `clq CHANGELOG.md`  
  validates the file.
- `clq -release CHANGELOG.md`  
  validates the file and further enforces that the most recent release is neither *[Unreleased]*
  nor has been *[YANKED]*. This validation is recommended before cutting a release or merging to main.
- `clq -query releases[0].version CHANGELOG.md`  
  validates the complete changelog and returns the version of the most recent release.

### Execution with Docker

A small minimal Docker image offers a simple no-installation executable. This image’s label is the release version,
with a secondary label ending in `-slim`, for example: `1.2.3-slim`.

A single changelog file can be validated with a simple `docker run -i denisa/clq < CHANGELOG.md`.

To operate on multiple files is more complex and we recommend either multiple individual
invocations, or the installation of native binaries.

Alternatively, the image is compatible with [Whalebrew](https://github.com/whalebrew/whalebrew).
After a one time installation `whalebrew install denisa/clq`, a one or more changelog files
can be validated with a simple `clq CHANGELOG.md` or `clq */CHANGELOG.md`.

The project also generates a 2nd Docker image whose label ends wih `-alpine`, for example: `1.2.3-alpine`.
This image, larger, is for use by the [clq-action](https://github.com/denisa/clq-action).

### GitHub Action

[clq-action](https://github.com/denisa/clq-action) documents how to integrate clq in a GitHub workflow.

## Grammar for supported Changelog

```text
CHANGELOG       = INTRODUCTION, RELEASES;
INTRODUCTION    = TITLE, { ? markdown paragraph ? };
TITLE           = "# ", ? inline content ?, LINE-ENDING;
RELEASES        = [ UNRELEASED ], { RELEASED | YANKED };
UNRELEASED      = UNRELEASED-HEAD, { CHANGES };
RELEASED        = RELEASED-HEAD, { CHANGES };
YANKED          = YANKED-HEAD, { CHANGES };
UNRELEASED-HEAD = "## [Unreleased]", LINE-ENDING;
RELEASED-HEAD   = "## [", SEMVER, "] - ", ISO-DATE, [ " [YANKED]" ], [ LABEL ], LINE-ENDING;
LABEL           = ? inline content, but not "[YANKED]" ?
CHANGES         = CHANGE-KIND, { CHANGE-DESC };
CHANGE-KIND     = "### ", ( "Added" | "Changed" | "Deprecated" | "Removed" | "Fixed" | "Security" ), LINE-ENDING;
CHANGE-DESC     = "- ", ? inline content ?, LINE-ENDING;
SEMVER          = ? see https://semver.org ?;
ISO-DATE        = YEAR, "-", MONTH, "-" DAY;
YEAR            = DIGIT, DIGIT, DIGIT, DIGIT;
MONTH           = DIGIT, DIGIT;
DAY             = DIGIT, DIGIT;
LINE-ENDING     = "U+000A" | "U+000D" | "U+000DU+000A";
DIGIT           = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";
```

Note:

- The most recent version comes first.
- The only deviation from the official spec is the optional LABEL on a release entry.
  The label is convenient for teams that want to highlight individual releases by naming
  them. Such might be the case for releases cut, for example, for a quarterly demo:
  `## [1.5.2] - 2019.10.02 Espelho`

## Validation

clq supports two validation modes, *feature* and *release*. The *feature* mode is best used for feature branches:
development work is in progress and the only requirement is for the changelog to be valid but not necessarily current.
On the opposite, the *release* mode applies to release branches and, therefore, pull-requests.

By default, clq operates in the *feature* mode. In that mode, clq validates that the changelog file conforms to the
grammar. It further validates that the releases are sorted chronologically from most recent to oldest,
that the versions numbers are properly decreasing and that the version change between any two versions
is justified by the change kinds present, according to a mapping from the type of change to the type of version change.

By default, the rules are:

- *major* release trigger:
  - `Added` for new features.
  - `Removed` for now removed features.
- *minor* release trigger:
  - `Changed` for changes in existing functionality.
  - `Deprecated` for soon-to-be removed features.
- *bugfix* release trigger:
  - `Fixed` for any bugfixes.
  - `Security` in case of vulnerabilities.

The `changeMap` option lets these rules be customized with a simple JSON file. In this example,
an *Added* section only triggers a minor version change:

```json
[
  {
    "name": "Added",
    "increment": "minor"
  },
  {
    "name": "Changed",
    "increment": "minor"
  },
  {
    "name": "Deprecated",
    "increment": "minor"
  },
  {
    "name": "Fixed",
    "increment": "patch"
  },
  {
    "name": "Removed",
    "increment": "major"
  },
  {
    "name": "Security",
    "increment": "patch"
  }
]
```

clq is generally lenient with the spaces, accepting them between square brackets for example.

When the *release* mode is activated (with the `-release` option), clq further validates
that the first entry in the changelog is an actual release entry.

*Note* that prereleases might or might not be supported at this time.  
![P’têt ben… P’têt pas… J’peux pas dire…](https://lestribulationsdunfrancophoneenfrancophonie.files.wordpress.com/2017/02/http-www-etaletaculture-frwp-contentuploads201512une-reponse-de-normands.jpg?w=317&h=269)  
(Astérix & Obélix, *Le tour de Gaule d’Astérix*, 1953)

## Emoji

The `changeMap` option further lets emoji be assigned to the change kinds with the optional `emoji` attribute in the
JSON file. The example extend on the previous one and define emoji for each change kinds:

```json
[
  {
    "name": "Added",
    "increment": "minor",
    "emoji": "✨"
  },
  {
    "name": "Changed",
    "increment": "major",
    "emoji": "💥"
  },
  {
    "name": "Deprecated",
    "increment": "minor",
    "emoji": "👎"
  },
  {
    "name": "Fixed",
    "increment": "patch",
    "emoji": "🐛"
  },
  {
    "name": "Removed",
    "increment": "major",
    "emoji": "🗑️"
  },
  {
    "name": "Security",
    "increment": "patch",
    "emoji": "🔒"
  }
]
```

## Experimental Extension to the standard

It is possible to use the change map file to define other change kinds, be they translation of the standard one, or new ones.
It is also possible to assign change kinds to the 'build' SemVer though those cannot be the single contribution of a release, as they would not increment the release number.
Please see [a French translation](docs/changemap/changedIsMajorWithEmoji_fr.json) and [using build for documentation](docs/changemap/withDocumentation.json).

## Query Expression Language

A query is a sequence of *query elements* leading through the structure of the changelog to the desired field.
The first query element is always a field from the changelog.

```text
QUERY            = ( SIMPLE_QUERY | COMPLEX_QUERY );
SIMPLE_QUERY     = { ARRAY_FIELD, "." }, FIELD;
COMPLEX_QUERY    = { ARRAY_FIELD, "." }, ARRAY_FIELD, ["/"];
ARRAY_FIELD      = FIELD, "[", [SELECTOR], "]";
FIELD            = ? see the Document Model section below ?;
SELECTOR         = DIGIT+;
DIGIT            = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";
```

A *simple* query returns the value of a single field. It is not formatted.

A *complex* query returns all the values of the selected object.
The object is formatted according to the value of the `-output` option.
If the query ends with a "/", it also returns the child elements.
If the selector is missing, the query returns a collection of objects.

For the sample changelog

```Markdown
# Change log

## [Unreleased]

### Added

- waldo
- fred

## [1.0.0] - 2020-06-20

### Removed

- foo
- bar
```

- `releases[1].version`  
  -> `1.0.0`
- `releases[1]`  
  -> `{"version":"1.0.0", "date":"2020-06-20"}`
- `releases[0].changes[]`  
  -> `[{"title":"Added"}]`
- `releases[0].changes[]/`  
  -> `[{"title":"Added", "descriptions":["waldo", "fred"]}]`

### Document Model

#### changelog

- *releases[]* all the releases defined in the changelog.  
  releases can be indexed, starting at 0, to access a single release.
- *title* the title of the changelog

#### release

- *changes[]* all the changes for that release.  
  changes cannot be indexed.
- *date* the release date, blank if it has not yet been released
- *label* the optional release label
- *status* one of *prereleased*, *released*, *unreleased* and *yanked*.
- *title* the version, date and optional label
- *version* the release version

#### change

- *descriptions[]* all the change descriptions;  
  descriptions cannot be indexed.
- *title*, the change kind.

## Reference

- [keep a changelog](https://keepachangelog.com/en/1.0.0/)
- [Semantic Versioning](https://semver.org)
- [CommonMark Spec](https://spec.commonmark.org/0.29/)
