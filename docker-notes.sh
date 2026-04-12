# notes on building, running, and getting size of images

# build
docker build . -t go

# run
docker run --rm go 1000000

# get image size
docker images | grep go | rev | cut -d " " -f 1 | rev
# => WARNING: This output is designed for human readability. For machine-readable output, please use --format.

docker image ls --filter "reference=*go*" --format "{{.Size}}"
docker image ls go --format "{{.Size}}"
docker image ls go:latest --format "{{.Size}}"

docker system df

docker image ls --filter "reference=*bday/*"

# make tabs 4 spaces in git CLI
git config --global core.pager 'less -x1,5'

docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' | sort > images.txt
cut -d' ' -f1 images.txt > tags-to-delete.txt
docker rmi $(cat tags-to-delte.txt)
cut -d' ' -f2 images.txt > images-to-delete.txt
docker image rm $(cat images-to-delete.txt)

sqlite3 app.db < schema.sql
sqlite3 app.db -header -box "select * from sessions"

# AI says 3 ways to insert 100 records:
# individual calls (no transaction) - slowest
# loop of calls inside transaction - 100x faster
# one sql statement - 2x faster than transaction
# ^ NEED TO TEST THIS

# 60/hr
curl -i -H "Accept: application/vnd.github+json" \
    https://api.github.com/repos/golang/go \
    > ghub-golang.json

# 5000/hr
curl -i -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $ghtoken" \
    https://api.github.com/repos/golang/go \
    > ghub-golang-authed.json

curl -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer $ghtoken" \
    https://api.github.com/repos/golang/go \
    | jq '.stargazers_count'

sqlite3 app.db ".schema --indent" > schema.sql
