FROM plugins/base:multiarch

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="ggz" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

EXPOSE 8080 8081

ADD release/linux/amd64/ggz /bin/

HEALTHCHECK --start-period=2s --interval=10s --timeout=5s \
  CMD ["/bin/ggz", "ping"]

ENTRYPOINT ["/bin/ggz"]
CMD ["server"]
