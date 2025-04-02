---
search:
  boost: 2
---

You can enable suggestions for mistyped flags or commands by setting `Suggest: true` on your `cli.App`. When enabled, if a user enters an unrecognized flag or command, the application will suggest the closest match if one is found based on Levenshtein distance.
