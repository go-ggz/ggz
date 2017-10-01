FROM plugins/base:multiarch
MAINTAINER Bo-Yi Wi <appleboy.tw@gmail.com>

EXPOSE 3003

ADD bin/api /

ENTRYPOINT ["/api"]
CMD ["server"]
