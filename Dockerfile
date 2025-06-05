FROM golang:1.24-alpine AS build

ARG TARGETOS
ARG TARGETARCH

ARG APP_RELEASE
ARG APP_RELEASE_DATE
ARG APP_VERSION
ARG CI_COMMIT_SHORT_SHA

WORKDIR /repo

COPY . .

ENV PACKAGE_NAME="github.com/matzefriedrich/containerssh-authserver"
ENV LDFLAGS="-X ${PACKAGE_NAME}/internal.CommitSha=${CI_COMMIT_SHORT_SHA} -X ${PACKAGE_NAME}/internal.Version=${APP_VERSION} -X ${PACKAGE_NAME}/internal.ReleaseDate=${APP_RELEASE_DATE} -X ${PACKAGE_NAME}/internal.ReleaseName=${APP_RELEASE}"

RUN mkdir -p build/ && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -ldflags "${LDFLAGS}" -a -o build/authserver cmd/authserver/main.go


FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=build /repo/build .

ENTRYPOINT [ "./authserver" ]