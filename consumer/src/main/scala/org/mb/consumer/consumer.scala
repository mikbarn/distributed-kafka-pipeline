package org.mb.consumer

//import scala.collection.JavaConverters._
import scala.jdk.CollectionConverters._
import org.apache.kafka.clients.consumer.KafkaConsumer
import org.apache.kafka.common.serialization.StringDeserializer
import java.util.Properties
import java.util.Collection
import java.time.Duration
import java.util.concurrent.atomic.{AtomicBoolean}

import collection.mutable._



object KafkaExample {
    val startTime = System.currentTimeMillis()
    val delta = () => System.currentTimeMillis() - startTime

    def main(args: Array[String]): Unit = {
        val properties = new Properties()
        var unsubscribed = new AtomicBoolean(true);
        var count: Int = 0
        properties.put("bootstrap.servers", sys.env.get("BOOTSTRAP_SERVER"))
        properties.put("group.id", "grp1")
        properties.put("key.deserializer", classOf[StringDeserializer])
        properties.put("value.deserializer", classOf[StringDeserializer])
        properties.put("auto.commit.interval.ms", "100")

        val kafkaConsumer = new KafkaConsumer[String, String](properties)

        this.synchronized {
            kafkaConsumer.subscribe(ListBuffer("t1").asJava)
        }
        unsubscribed.set(false)
        def shutdown() = {
            this.synchronized {
                if (!unsubscribed.get()) {
                    println("Unsubscribing...")
                    kafkaConsumer.unsubscribe()
                    unsubscribed.set(true)
                } else {
                    println("Already unsubscribed!")
                }
            }
        }

        sys.addShutdownHook({shutdown()})

        var total: Int = 0
        try {
            while (count < 100 && delta() < 1000 * 60 * 10)  {
                println(s"${delta()} ms have elapsed!")
                //val dur = Duration.ofMillis(10)
                val dur = Duration.ofSeconds(30)
                val consRecs = this.synchronized {
                    kafkaConsumer.poll(dur)
                }
                if (consRecs.isEmpty()) {
                    println("No messages received! Must have timed out...")

                } else {
                    count += 1
                    println(s"Partitions found: ")
                    for (p <- consRecs.partitions.asScala) {
                        println(p)
                    }
                    val resIt = consRecs.iterator().asScala
                    while (resIt.hasNext)  {
                        total += 1
                        //println(resIt.next)
                        val n = resIt.next()
                        println("----------------------")
                        println("Received this record: ")
                        println(s"Topic: ${n.topic()}")
                        println(s"Part: ${n.partition()}")
                        println(s"Key: ${n.key()}")
                        println(s"Offset: ${n.offset()}")
                        println(s"TS: ${n.timestamp()}")
                        println(s"Val: ${n.value()}")
                    }
                }

            }
            println(s"Received $count times! Total messages: $total")
            kafkaConsumer.commitSync()
        }
        finally {
            shutdown()
        }
    }
}


