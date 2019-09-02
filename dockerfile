FROM debian:stretch-slim
# RUN apk --no-cache add curl ca-certificates bash
ADD apitest /bin/apitest
CMD ["/bin/apitest"]