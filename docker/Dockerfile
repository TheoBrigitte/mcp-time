FROM scratch
ARG TARGETOS
ARG TARGETARCH

ADD build/mcp-time.${TARGETOS}-${TARGETARCH} /app/mcp-time

ENTRYPOINT ["/app/mcp-time"]
