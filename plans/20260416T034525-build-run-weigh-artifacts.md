## Theme
- always operate line by line and make an artifact

## Convert results files to JSONL
- we should convert results files to jsonl
- at least the tempfile
- what about the results and previous results files?
- yeah, let's just switch them all over

## Build script make artifact
- build script doesn't make an artifact
- we should have it build a file line by line of image names it made

## Get docker image size one by one
- we currently just list all docker image sizes in one go
- can we loop over the images and get the size per image?
- we will still have to do math in go, but we avoid the goofy string math
- this also simplifies the pre-main-post setup
