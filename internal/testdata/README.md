
Sample files were obtained from [here](http://techslides.com/sample-files-for-development).

The test asset file `bindata.go` was produced using [go-bindata](https://github.com/go-bindata/go-bindata) with the command: 

`go-bindata -pkg testdata ./original/... ./tagged/...`

Samples are tagged with the following test metadata (where supported by the tag format):

| Field            | Value            | ID3 v1 | ID3 v1.1 | ID3 v2.x |
| ---------------- | ----------------:|:------:|:--------:|:--------:|
| **Artist**       | Test Artist      |   x    |    x     |    x     |
| **Track Title**  | Test Title       |   x    |    x     |    x     |
| **Album Title**  | Test Album       |   x    |    x     |    x     |
| **Track Number** | 3                |        |    x     |    x     |
| **Genre**        | Jazz             |   x    |    x     |    x     |
| **Year**         | 2000             |   x    |    x     |    x     |
| **Comment**      | Test Comment     |   x    |    x     |    x     |
| **Album Artist** | Test AlbumArtist |        |          |    x     |
| **Composer**     | Test Composer    |        |          |    x     |
| **Disc Number**  | 2                |        |          |    x     |
| **Track Count**  | 6                |        |          |    x     |