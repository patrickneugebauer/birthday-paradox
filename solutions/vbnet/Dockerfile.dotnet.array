FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine
WORKDIR /app

COPY loops.vbproj .
COPY Program.vb .

RUN dotnet build

ENTRYPOINT ["dotnet", "bin/Debug/net10.0/loops.dll"]
