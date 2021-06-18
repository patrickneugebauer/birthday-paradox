FROM mcr.microsoft.com/dotnet/core/sdk:3.1

WORKDIR /root

ADD loops.csproj .
ADD Program.cs .

RUN dotnet build

ENTRYPOINT ["dotnet", "bin/Debug/netcoreapp3.1/loops.dll"]
