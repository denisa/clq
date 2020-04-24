# clq
Changelog query tool

# Supported grammar
CHANGELOG  = TITLE [ UNRELEASED ] { RELEASED | YANKED  }
TITLE           = "#" inline content LINE_ENDING { markdown paragraph }
UNRELEASED      = UNRELEASED_HEAD { CHANGES }
RELEASED        = RELEASED_HEAD { CHANGES }
YANKED          = YANKED_HEAD { CHANGES }
UNRELEASED_HEAD = "## [Unreleased]" LINE_ENDING
RELEASED_HEAD   = "## [" SEMVER "] - " ISO_DATE [ LABEL ] LINE_ENDING
YANKED_HEAD     = "## " SEMVER " - " ISO_DATE "[YANKED]" LINE_ENDING
CHANGES         = [ ADDED ] [ CHANGED ] [ DEPRECATED ] [ REMOVED ] [ FIXED ] [ SECURITY ]
ADDED           = "### Added" LINE_ENDING
CHANGED         = "### Changed" LINE_ENDING
DEPRECATED      = "### Deprecated" LINE_ENDING
REMOVED         = "### Removed" LINE_ENDING
FIXED           = "### Fixed" LINE_ENDING
SECURITY        = "### Security" LINE_ENDING
CHANGE_DESC     = "- " inline content LINE_ENDING
ISO_DATE        = YYYY "-" MM "-" DD
LINE_ENDING     = "U+000A" | "U+000D" | "U+000DU+000A"
Note:
- The latest version comes first.
- Definition of SEMVER at https://semver.org

Reference:
- https://keepachangelog.com/en/1.0.0/
- https://semver.org
- https://spec.commonmark.org/0.29/
