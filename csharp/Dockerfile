FROM mcr.microsoft.com/dotnet/core/sdk:2.1

WORKDIR /root

ADD loops.csproj .
ADD Program.cs .

RUN dotnet build

ENTRYPOINT ["dotnet", "bin/Debug/netcoreapp2.1/loops.dll"]
