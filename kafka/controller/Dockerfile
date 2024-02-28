FROM ibmjava:8-jre

RUN groupadd kafkagrp
RUN useradd -g kafkagrp -m -d /home/kafkauser -s /bin/bash kafkauser
RUN mkdir /opt/kafka

RUN chown -R kafkauser:kafkagrp /opt/kafka/

WORKDIR /opt/kafka/
USER kafkauser


ARG KCID=FroGkXidSBuRYv_ndR36-Q
ENV KAFKA_CLUSTER_ID=${KCID}


EXPOSE 9092 9093
ENTRYPOINT ["tail", "-f", "/dev/null"]

