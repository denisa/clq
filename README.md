# clq — Changelog validation and query tool

[![version](https://img.shields.io/github/v/release/denisa/clq?include_prereleases&sort=semver)](https://github.com/denisa/clq/releases)
[![semantic versioning](https://img.shields.io/badge/semantic%20versioning-2.0.0-informational)](https://semver.org/spec/v2.0.0.html)

[![test](https://github.com/denisa/clq/workflows/test/badge.svg)](https://github.com/denisa/clq/actions?query=workflow%3Atest+branch%3Amaster)
[![docker build](https://img.shields.io/docker/cloud/build/denisa/clq)](https://hub.docker.com/repository/docker/denisa/clq/builds)
[![coverage status](https://coveralls.io/repos/github/denisa/clq/badge.svg?branch=master)](https://coveralls.io/github/denisa/clq?branch=master)

# usage

clq always validates the complete changelog, stopping at its first error. If a query is given, clq then queries the changelog and returns the query result. clq handles standard input — when no arguments are present or an argument is "-" — or any number of files.

clq exits with a status of 0 if all files are valid, and with a non-zero status if any file fails to validate. It writes to standard output the result of the query if a query was given.

clq writes validation error to standard error.

When processing multiple files, clq prefixes every line on standard out and standard error with the file name.

    Usage: clq { options } <path to changelog.md>

    Options are:
      -output format
          the format to apply to the result of a (complex) query. Supports json and md (markdown); default to json
      -query string
        	A query to extract information out of the change log
      -release
        	Enable release-mode validation
      -with-filename
        	Always print filename headers with output lines

Example:

- `clq CHANGELOG.md`  
    validates the file.
- `clq -release CHANGELOG.md`  
    validates the file and further enforces that the most recent release is neither _[Unreleased]_ nor has been _[YANKED]_. This validation is recommended before cutting a release or merging to master.
- `clq -query releases[0].version CHANGELOG.md`  
    validates the complete changelog and returns the version of the most recent release.

# Execution with Docker

A small docker image offers a simple no-installation executable.

A single changelog file can be validated with a simple `docker run -i denisa/clq < CHANGELOG.md`.

To operate on multiple files is more complex and we recommend either multiple individual invocations, or the installation of native binaries.

# Grammar for supported Changelog
```
CHANGELOG       = INTRODUCTION, RELEASES;
INTRODUCTION    = TITLE, { ? markdown paragraph ? };
TITLE           = "# ", ? inline content ?, LINE-ENDING;
RELEASES        = [ UNRELEASED ], { RELEASED | YANKED };
UNRELEASED      = UNRELEASED-HEAD, { CHANGES };
RELEASED        = RELEASED-HEAD, { CHANGES };
YANKED          = YANKED-HEAD, { CHANGES };
UNRELEASED-HEAD = "## [Unreleased]", LINE-ENDING;
RELEASED-HEAD   = "## [", SEMVER, "] - ", ISO-DATE, [ LABEL ], LINE-ENDING;
LABEL           = ? inline content, but not "[YANKED]" ?
YANKED-HEAD     = "## ", SEMVER, " - ", ISO-DATE, " [YANKED]", LINE-ENDING;
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

# Validation

clq validates that the changelog file conforms to the grammar. It further validates that the releases are sorted chronologically from most recent to oldest, that the versions numbers are properly decreasing and that the version change between any two versions is justified by the change kinds present, according to folowing rules:

- _major_ release trigger:
    - `Added` for new features.
    - `Removed` for now removed features.
- _minor_ release trigger:
    - `Changed` for changes in existing functionality.
    -  `Deprecated` for soon-to-be removed features.
- _bug-fix_ release trigger:
    - `Fixed` for any bug fixes.
    - `Security` in case of vulnerabilities.

clq is generally lenient with the spaces, accepting them between square brackets for example.

_Note_ that pre-releases might or might not be supported at this time.  
![P’têt ben… P’têt pas… J’peux pas dire…](https://lestribulationsdunfrancophoneenfrancophonie.files.wordpress.com/2017/02/http-www-etaletaculture-frwp-contentuploads201512une-reponse-de-normands.jpg?w=317&h=269)  
(Astérix & Obélix, _Le tour de Gaule_, 1953)

# Query Expression Language

A query is a sequence of _query elements_ leading through the structure of the changelog to the desired field.
The first query element is always a field from the changelog.

```
QUERY            = ( SIMPLE_QUERY | COMPLEX_QUERY );
SIMPLE_QUERY     = { ARRAY_FIELD, "." }, FIELD;
COMPLEX_QUERY    = { ARRAY_FIELD, "." }, ARRAY_FIELD, ["/"];
ARRAY_FIELD      = FIELD, "[", [SELECTOR], "]";
FIELD            = ? see the Document Model section below ?;
SELECTOR         = DIGIT+;
DIGIT            = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";
```
A _simple_ query returns the value of a single field. It is not formatted.

A _complex_ query returns all the values of the selected object. The object is formatted accoring to the value of `-output` option. If the query ends with a "/", it also returns the values child elements. If the selector is missing, the query returns a collection of objects.

For the sample changelog
```
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

## Document Model
### changelog
- _releases[]_ all the releases defined in the changelog.  
  releases can be indexed, starting at 0, to access a single release.
- _title_ the title of the changelog
### release
- _changes[]_ all the changes for that releases.  
  changes cannot be indexed.
- _date_ the release date, blank if it has not yet been released
- _label_ the optional release label
- _status_ one of _prereleased_, _released_, _unreleased_ and _yanked_.
- _title_ the version, date and optional label
- _version_ the release version
### change
- _descriptions[]_ all the change descriptions.  
  descriptions cannot be indexed.
- _title_, the change kind

# Reference:

-   <https://keepachangelog.com/en/1.0.0/>
-   <https://semver.org>
-   <https://spec.commonmark.org/0.29/>
