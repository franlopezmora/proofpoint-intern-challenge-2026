# Cleaning Report

- Total number of input records: <to-be-generated>
- Total number of output records: <to-be-generated>
- Number of discarded entries: <to-be-generated>
- Number of corrected entries: <to-be-generated>
- Number of duplicates detected: <to-be-generated>
 
_Template note: these values are generated at runtime by the CLI._

## Deduplication strategy
Records are compared using normalized values (trimmed, whitespace-collapsed, and case-normalized text). Two records are considered duplicates if they match any of these rules:

1. `(SeriesName, SeasonNumber, EpisodeNumber)`
2. `(SeriesName, 0, EpisodeNumber, EpisodeTitle)`
3. `(SeriesName, SeasonNumber, 0, EpisodeTitle)`

When multiple records refer to the same episode, the program keeps the best one using the following priority:
1. valid AirDate over `Unknown`
2. known EpisodeTitle over `Untitled Episode`
3. known SeasonNumber and EpisodeNumber over fallback `0`
4. first appearance in the input file if all previous criteria are tied

## Corrected entries policy
A record is counted as corrected when at least one field had to be normalized or replaced during cleaning. This includes trimming or collapsing whitespace, normalizing text casing for comparison/output consistency, replacing invalid or missing numeric values with `0`, replacing a missing title with `Untitled Episode`, and replacing an invalid or missing air date with `Unknown`.
