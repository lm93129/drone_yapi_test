FROM alpine:3.4
RUN apk --no-cache add curl ca-certificates bash
ADD apitest /bin/apitest
ENTRYPOINT ["/bin/apitest"]