WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev migrate

# dependencies

COPY ["app/go.mod" , "./"]

# build