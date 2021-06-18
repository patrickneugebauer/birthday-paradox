FROM mcr.microsoft.com/dotnet/core/sdk:3.1

WORKDIR /root

ADD loops.fsproj .
ADD Program.fs .

RUN dotnet build

ENTRYPOINT ["dotnet", "bin/Debug/netcoreapp3.1/loops.dll"]
