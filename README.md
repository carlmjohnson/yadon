# yadon [![GoDoc](https://godoc.org/github.com/carlmjohnson/yadon?status.svg)](https://godoc.org/github.com/carlmjohnson/yadon)
How slow can you GET?

[![Yadon video](https://img.youtube.com/vi/Ce5mRvkAePU/0.jpg)](https://www.youtube.com/watch?v=Ce5mRvkAePU)

Yadon is a tool to request a site as slowly as possible. The purpose of the tool is to test whether the site might be vulnerable to a [slowloris attack][]. If you test your server and find that it's possible for a slow connection to draw out a download of a few kilobytes for minutes or hours, you should decrease your server idle timeouts.

[slowloris attack]: https://en.wikipedia.org/wiki/Slowloris_(computer_security)
