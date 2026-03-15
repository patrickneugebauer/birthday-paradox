FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine

WORKDIR /root

ADD loops.fsproj .
ADD Program.fs .

RUN dotnet build

ENTRYPOINT ["dotnet", "bin/Debug/net10.0/loops.dll"]
