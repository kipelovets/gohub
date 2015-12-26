FROM golang

RUN mkdir /gohub
COPY ./* /gohub/

RUN cd /gohub && go build -v && cp gohub /usr/bin/gohub

CMD /usr/bin/gohub --log=- --port=7654 --config=/gohub/example.json