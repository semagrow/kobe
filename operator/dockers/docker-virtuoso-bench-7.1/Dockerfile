FROM ubuntu:12.04

MAINTAINER Yiannis Mouchakis <gmouchakis@iit.demokritos.gr>

# Install Virtuoso prerequisites
RUN apt-get update \
        && apt-get install -y curl build-essential debhelper autotools-dev autoconf automake unzip wget net-tools git libtool flex bison gperf gawk m4 libssl-dev libreadline-dev openssl

# Virtuoso 7.1 commit
ENV VIRTUOSO_TAG v7.1.0

RUN git clone https://github.com/openlink/virtuoso-opensource.git \
        && cd virtuoso-opensource \
        && git checkout tags/${VIRTUOSO_TAG} \
        && ./autogen.sh \
        && ./configure --with-readline --program-transform-name="s/isql/isql-v/" \
        && make && make install \
        && ln -s /usr/local/virtuoso-opensource/var/lib/virtuoso/ /var/lib/virtuoso \
	&& ln -s /var/lib/virtuoso/db /data \
        && cd .. \
        && rm -r /virtuoso-opensource

# Add Virtuoso bin to the PATH
ENV PATH /usr/local/virtuoso-opensource/bin/:$PATH

# Add Virtuoso config
ADD virtuoso.ini /virtuoso.ini

# Add startup script
ADD virtuoso.sh /virtuoso.sh

WORKDIR /data

EXPOSE 8890

CMD ["/bin/bash", "/virtuoso.sh"]
