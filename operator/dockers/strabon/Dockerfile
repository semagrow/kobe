#######################################
# Dockerfile for Strabon              #
#                                     #
# based on bde2020/sc7pilotbase:1.0.0 #
# java 1.8.0_151                      #
# maven 3.0.5                         #
# git 1.9.1                           #
# tomcat 8.0.23                       #
#######################################
FROM bde2020/sc7pilotbase:1.0.0

RUN apt-get update \
 && apt-get install --assume-yes mercurial wget \
 && apt-get --assume-yes --purge autoremove \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

RUN hg clone 'http://hg.strabon.di.uoa.gr/Strabon/' \
 && cd Strabon \
 && mvn clean package

ENV STRABON_HOME=/tomcat/webapps/strabon
RUN mkdir /resources \
 && unzip -d $STRABON_HOME /Strabon/endpoint/target/strabon-endpoint-3.3.2-SNAPSHOT.war

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/bin/bash", "/entrypoint.sh"]

ADD runstrabon.sh /runstrabon.sh
ADD loaddata.sh /loaddata.sh
CMD ["/bin/bash", "/runstrabon.sh"]

