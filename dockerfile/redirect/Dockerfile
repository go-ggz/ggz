FROM plugins/base:multiarch

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="ggz" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

EXPOSE 8080

ADD release/linux/amd64/ggz-redirect /bin/

HEALTHCHECK --start-period=2s --interval=10s --timeout=5s \
  CMD ["/bin/ggz-redirect", "ping"]

ENTRYPOINT ["/bin/ggz-redirect"]
CMD ["server"]
