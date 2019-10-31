FROM swift

WORKDIR /root

COPY .gitignore .
COPY Package.swift .
COPY README.md .
COPY Sources/ Sources/
COPY Tests/ Tests/

RUN swift build -c release

ENTRYPOINT ["./.build/release/swift"]
