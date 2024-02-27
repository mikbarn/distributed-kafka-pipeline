ThisBuild / scalaVersion := "2.13.12"

lazy val root = (project in file("."))
  .settings(
    name := "KafkaConsumer",
    assembly / mainClass := Some("org.mb.consumer.KafkaExample"),
    libraryDependencies += "org.apache.kafka" % "kafka-clients" % "3.6.1"
  )


  ThisBuild / assemblyMergeStrategy := {
//   case PathList("javax", "servlet", xs @ _*)         => MergeStrategy.first
//   case PathList(ps @ _*) if ps.last endsWith ".html" => MergeStrategy.first
  //case "application.conf"                            => MergeStrategy.concat
  case "unwanted.txt"                                => MergeStrategy.discard
  case x =>
    val oldStrategy = (ThisBuild / assemblyMergeStrategy).value
    oldStrategy(x)
}