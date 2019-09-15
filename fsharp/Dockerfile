FROM mcr.microsoft.com/dotnet/core/sdk:2.1

WORKDIR /root

ADD loops.fsproj .
ADD Program.fs .

RUN dotnet build

ENTRYPOINT ["dotnet", "bin/Debug/netcoreapp2.1/loops.dll"]
