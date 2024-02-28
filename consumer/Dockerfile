FROM ubuntu:20.04

# RUN apt update
# RUN apt install -y default-jre
# RUN apt install -y dnsutils
# RUN apt install -y netcat
# RUN apt install -y iputils-ping
# RUN apt install -y net-tools
# RUN apt install iproute2 -y

RUN apt update && \ 
    apt install -y  \
    default-jre \
    dnsutils \
    && rm -rf /var/lib/apt/lists/* \
    && apt clean

RUN groupadd kafkagrp
RUN useradd -g kafkagrp -m -d /home/kafkauser -s /bin/bash kafkauser
RUN mkdir /opt/app
COPY target/scala-2.13/KafkaConsumer-assembly-0.1.0-SNAPSHOT.jar /opt/app/
RUN chown -R kafkauser:kafkagrp /opt/app/

WORKDIR /opt/app
USER kafkauser

EXPOSE 9092
ENTRYPOINT ["java", "-jar",  "KafkaConsumer-assembly-0.1.0-SNAPSHOT.jar"]

