# Cleaning Report

- Total number of input records: 6
- Total number of output records: 2
- Number of discarded entries: 3
- Number of corrected entries: 5
- Number of duplicates detected: 1

## Deduplication strategy
Records are compared using normalized keys. A primary key uses (series, season, episode). Secondary keys are used only when one numeric field is missing: (series, 0, episode, title) when season is 0, and (series, season, 0, title) when episode is 0. For each duplicate cluster, the kept record is selected by priority: known AirDate, known EpisodeTitle, both season and episode known, then first appearance in input order.

## Corrected entries policy
An entry is counted as corrected when at least one field was normalized or defaulted (trim/collapse spaces, case normalization, number fallback to 0, title fallback to "Untitled Episode", or air date fallback to "Unknown"/normalized ISO date).
