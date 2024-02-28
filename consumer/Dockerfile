FROM ibmjava:8-jre

RUN groupadd kafkagrp
RUN useradd -g kafkagrp -m -d /home/kafkauser -s /bin/bash kafkauser
RUN mkdir /opt/app
COPY target/scala-2.13/KafkaConsumer-assembly-0.1.0-SNAPSHOT.jar /opt/app/
RUN chown -R kafkauser:kafkagrp /opt/app/

WORKDIR /opt/app
USER kafkauser

EXPOSE 9092
ENTRYPOINT ["java", "-jar",  "KafkaConsumer-assembly-0.1.0-SNAPSHOT.jar"]

