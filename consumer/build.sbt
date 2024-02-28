ThisBuild / scalaVersion := "2.13.12"


lazy val consumer = project
  .in(file("."))
  .settings(
    name := "KafkaConsumer",
    assembly / mainClass := Some("org.mb.consumer.KafkaExample"),
    libraryDependencies += "org.apache.kafka" % "kafka-clients" % "2.6.0",
    libraryDependencies += "com.fasterxml.jackson.core" % "jackson-databind" % "2.10.2"
  )


  ThisBuild / assemblyMergeStrategy := {
    //case PathList("javax", "servlet", xs @ _*)         => MergeStrategy.first
   // case PathList(ps @ _*) if ps.last endsWith ".html" => MergeStrategy.first
   // case "application.conf"                            => MergeStrategy.concat
    case "module-info.class"                                => MergeStrategy.discard
    case x =>
      val oldStrategy = (ThisBuild / assemblyMergeStrategy).value
      oldStrategy(x)
}