FROM microsoft/nanoserver:10.0.14393.1884

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="ggz-redirect" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

EXPOSE 8080

COPY release/ggz-redirect.exe C:/bin/ggz-redirect.exe

ENTRYPOINT [ "C:\\bin\\ggz-redirect.exe" ]
