# clq — Changelog validation and query tool

[![test](https://github.com/denisa/clq/workflows/test/badge.svg)](https://github.com/denisa/clq/actions?query=workflow%3Atest+branch%3Amaster)
[![Coverage Status](https://coveralls.io/repos/github/denisa/clq/badge.svg?branch=master)](https://coveralls.io/github/denisa/clq?branch=master)
![Semantic Versioning](https://img.shields.io/badge/Sematic%20Versioning-2.0.0-informational)

# Grammar for supported Changelog
```
CHANGELOG       = TITLE, [ UNRELEASED ], { RELEASED | YANKED };
TITLE           = "#", ? inline content ?, LINE-ENDING, { ? markdown paragraph ? };
UNRELEASED      = UNRELEASED-HEAD, { CHANGES };
RELEASED        = RELEASED-HEAD, { CHANGES };
YANKED          = YANKED-HEAD, { CHANGES };
UNRELEASED-HEAD = "## [Unreleased]", LINE-ENDING;
RELEASED-HEAD   = "## [", SEMVER, "] - ", ISO-DATE, [ LABEL ], LINE-ENDING;
YANKED-HEAD     = "## ", SEMVER, " - ", ISO-DATE, "[YANKED]", LINE-ENDING;
CHANGES         = CHANGE-KIND, { CHANGE-DESC };
CHANGE-KIND     = "### ", ( "Added" | "Changed" | "Deprecated" | "Removed" | "Fixed" | "Security" ), LINE-ENDING;
CHANGE-DESC     = "- ", ? inline content ?, LINE-ENDING;
ISO-DATE        = YEAR, "-", MONTH, "-" DAY;
YEAR            = DIGIT, DIGIT, DIGIT, DIGIT;
MONTH           = DIGIT, DIGIT;
DAY             = DIGIT, DIGIT;
LINE-ENDING     = "U+000A" | "U+000D" | "U+000DU+000A";
DIGIT           = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9";
```
Note:
- The latest version comes first.
- Definition of SEMVER at https://semver.org

Reference:
- https://keepachangelog.com/en/1.0.0/
- https://semver.org
- https://spec.commonmark.org/0.29/
