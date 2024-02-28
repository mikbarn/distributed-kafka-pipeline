FROM ubuntu:20.04

RUN apt update && \ 
    apt install -y  \
    default-jre \
    dnsutils \
    && rm -rf /var/lib/apt/lists/* \
    && apt clean

RUN groupadd kafkagrp
RUN useradd -g kafkagrp -m -d /home/kafkauser -s /bin/bash kafkauser
RUN mkdir /opt/kafka
# COPY dls/kafka_2.13-3.6.1.tgz /opt/kafka/
# COPY ./init_ctlr.sh /opt/kafka/

RUN chown -R kafkauser:kafkagrp /opt/kafka/

WORKDIR /opt/kafka/
USER kafkauser
#RUN tar -xzf kafka_2.13-3.6.1.tgz
# RUN rm  kafka_2.13-3.6.1.tgz

WORKDIR /opt/kafka/

ARG KCID=FroGkXidSBuRYv_ndR36-Q
ENV KAFKA_CLUSTER_ID=${KCID}


EXPOSE 9092 9093
ENTRYPOINT ["bash", "-c", "/opt/kafka/init_ctlr.sh"]

