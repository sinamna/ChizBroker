FROM golang
COPY ./build/server/broker /bin/broker
CMD /bin/broker