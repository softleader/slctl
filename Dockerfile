FROM softleader/alpine

COPY _build/slctl /usr/local/bin
COPY _build/echo/ /root/.sl/plugins/