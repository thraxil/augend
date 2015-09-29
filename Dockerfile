FROM scratch
MAINTAINER Anders Pearson <anders@columbia.edu>
COPY augend /
COPY media /media
COPY templates /templates
EXPOSE 8890
CMD ["/augend", "-config=/config.conf"]

