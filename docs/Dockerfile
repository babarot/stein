FROM alpine:3.10.2

RUN apk update && apk add --no-cache \
  bash \
  git \
  git-fast-import \
  openssh \
  python3 \
  python3-dev \
  curl \
  && python3 -m ensurepip \
  && rm -r /usr/lib/python*/ensurepip \
  && pip3 install --upgrade pip setuptools \
  && rm -r /root/.cache \
  && rm -rf /var/cache/apk/*

COPY requirements.txt /
RUN  pip install -U -r /requirements.txt

WORKDIR /docs

EXPOSE 3000

CMD ["mkdocs", "serve", "--dev-addr=0.0.0.0:3000", "--livereload"]
