FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine
WORKDIR /app

COPY loops.ps1 .

ENTRYPOINT ["pwsh", "loops.ps1"]
