# photo-organize

A tool for delivering photos to clients.

Assume we have a list of photos ordered by time, every group's first and last photo should contains a QR Code with same content as group id. This tool will read and group all photos, then upload them to OSS.

```
// group 1
1.1.jpg
1.2.jpg
1.3.jpg
1.4.jpg
// group 2
2.1.jpg
2.2.jpg
2.3.jpg
2.4.jpg
```