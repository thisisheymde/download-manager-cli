# stupid simple downloader

a cli tool written in go. uses goroutines for concurrent downloads. It breaks a file down to chunks and downloads them concurrently, making your downloads faster.

for a single file download 

`ddman -url "urlname" -dir "dirname" -connections num`

for multiple file download

`ddman -multiple "filename" -dir "dirname" -connections num`

P.S. 

filename should be in format like this

```
url1
url2
url3
and so on
```

default directory is Downloads

default number of connections is 4