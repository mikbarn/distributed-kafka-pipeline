FROM ubuntu:20.04

RUN apt update
RUN apt install -y default-jre
RUN apt install -y dnsutils
RUN apt install -y netcat
RUN apt install -y iputils-ping
RUN apt install -y net-tools
RUN apt install -y iproute2

RUN groupadd kafkagrp
RUN useradd -g kafkagrp -m -d /home/kafkauser -s /bin/bash kafkauser
RUN mkdir /opt/kafka
COPY dls/kafka_2.13-3.6.1.tgz /opt/kafka/
COPY ./init_ctlr.sh /opt/kafka/

RUN chown -R kafkauser:kafkagrp /opt/kafka/

WORKDIR /opt/kafka/
USER kafkauser
RUN tar -xzf kafka_2.13-3.6.1.tgz
WORKDIR /opt/kafka/kafka_2.13-3.6.1/
ARG KCID=FroGkXidSBuRYv_ndR36-Q
ENV KAFKA_CLUSTER_ID=${KCID}


EXPOSE 9092 9093
ENTRYPOINT ["bash", "-c", "/opt/kafka/init_ctlr.sh"]

