FROM plugins/base:multiarch

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>"

EXPOSE 8080 8081 80 443

ADD bin/ggz /

ENTRYPOINT ["/ggz"]
CMD ["server"]
