FROM openjdk:11

COPY ./apache-maven /usr/src/apache-maven

RUN mkdir /usr/src/resource
WORKDIR /usr/src/resource

CMD ["/usr/src/apache-maven/bin/mvn","-s","/usr/src/apache-maven/conf/settings.xml","deploy"]