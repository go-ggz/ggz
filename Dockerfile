FROM plugins/base:multiarch

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>"

EXPOSE 8080 8081

ADD bin/ggz /

HEALTHCHECK --start-period=2s --interval=10s --timeout=5s \
  CMD ["/ggz", "ping"]

ENTRYPOINT ["/ggz"]
CMD ["server"]
