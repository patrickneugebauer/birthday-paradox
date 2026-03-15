# run in one command
dotnet Program.cs -- 1000000
dotnet run --file Program.cs -- 1000000

# build and run with framework
dotnet build
dotnet bin/Debug/net10.0/loops.dll 1000000
./bin/Debug/net10.0/loops 1000000

# publish standalone and run
dotnet publish -c Release \
    -r linux-musl-x64 \
    --self-contained true \
    -p:PublishSingleFile=true \
    -p:PublishTrimmed=true \
    -o ./publish
./publish/loops 1000000

# in .NET 10 you don't need a class and main method to start the app

# To run Dockerfile with different name
# docker build -f Dockerfile.build . -t csharp-build
# docker build -f Dockerfile.publish . -t csharp-publish
# docker run -it --rm --entrypoint sh -t csharp-build
# du -sh $(ls -A)

# base image image 890M
# base alpine image 758M
# base alpine with build 758M
# base alpine with publish 1.11G

# base runtime image 120M
# base alpine runtime 11.4M
# base alpine with executable 26.3M
