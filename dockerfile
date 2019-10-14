FROM debian:stretch-slim
ADD apitest /bin/apitest
CMD ["/bin/apitest"]