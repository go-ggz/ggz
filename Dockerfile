FROM plugins/base:multiarch

LABEL maintainer="Bo-Yi Wi <appleboy.tw@gmail.com>"

EXPOSE 8080 8081

ADD bin/ggz /

ENTRYPOINT ["/ggz"]
CMD ["server"]
