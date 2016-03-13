FROM alpine

ADD twitchbot /bin/twitchbot

ENTRYPOINT ["/bin/twitchbot"]