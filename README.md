# clq — Changelog validation and query tool

[![test](https://github.com/denisa/clq/workflows/test/badge.svg)](https://github.com/denisa/clq/actions?query=workflow%3Atest+branch%3Amaster)
[![docker build](https://img.shields.io/docker/cloud/build/denisa/clq)](https://hub.docker.com/repository/docker/denisa/clq/builds)
[![coverage status](https://coveralls.io/repos/github/denisa/clq/badge.svg?branch=master)](https://coveralls.io/github/denisa/clq?branch=master)
[![semantic versioning](https://img.shields.io/badge/semantic%20versioning-2.0.0-informational)](https://semver.org/spec/v2.0.0.html)

# usage
clq always validates the complete changelog, stopping at its first error. If a query is given, clq then queries the changelog and returns the query result. clq handles standard input — when no arguments are present or an argument is "-" — or any number of files.

clq exits with a status of 0 if all files are valid, and with a non-zero status if any file fails to validate. It writes to standard output the result of the query if a query was given.

clq writes validation error to standard error.

When processing multiple files, clq prefixes every line on standard out and standard error with the file name.

```
Usage: clq { options } <path to changelog.md>

Options are:
  -query string
    	A query to extract information out of the change log
  -release
    	Enable release-mode validation
  -with-filename
    	Always print filename headers with output lines
```

Example:
- `clq CHANGELOG.md`\
validates the file.
- `clq -release CHANGELOG.md`\
validates the file and further enforces that the most recent release is neither _[Unreleased]_ nor has been _[YANKED]_. This validation is recommended before cutting a release or merging to master.
- `clq -query releases[0].version CHANGELOG.md`\
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
   - `Deprecated` for soon-to-be removed features.
- _bug-fix_ release trigger:
   - `Fixed` for any bug fixes.
   - `Security` in case of vulnerabilities.

clq is generally lenient with the spaces, accepting them between square brackets for example.

_Note_ that pre-releases might or might not be supported at this time.\
![P’têt ben… P’têt pas… J’peux pas dire…](https://lestribulationsdunfrancophoneenfrancophonie.files.wordpress.com/2017/02/http-www-etaletaculture-frwp-contentuploads201512une-reponse-de-normands.jpg?w=317&h=269)\
(Astérix & Obélix, _Le tour de Gaule_, 1953)

# Query Expression Language
```
QUERY            = "title" | ( "releases[", RELEASE_SELECTOR, "]", [ "." ( "date", "label", "version", "status")] );
RELEASE_SELECTOR = DIGIT+;
DIGIT            = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";
```
The query expression is a sequence of _query elements_.

The first query element accesses the overall _changelog_ and provides accessor for its _title_ and an array of _releases_. Returning all releases is not supported, and a _release selector_ identifies one of the releases. Currently, the release selector is the (0-based) index of the release, with the most recent entry having the 0th index. clq returns an empty result if the index is bigger than the release count.

The second query element, if present, accesses fields of the release. There are accessors for the _version_, the _date_, the _label_, and the release _status_. The release status is one of _released_, _unreleased_ and _yanked_. clq returns a json structure with all the defined fields of the selected release if the second query element is missing.


# Reference:
- https://keepachangelog.com/en/1.0.0/
- https://semver.org
- https://spec.commonmark.org/0.29/
