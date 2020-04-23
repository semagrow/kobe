appender( "PROCFLOW", ConsoleAppender ) {
  encoder( PatternLayoutEncoder) {
    pattern = "%.-1level - %date{ISO8601} - [%4.4thread] - %10.-10logger{0} - %.12X{uuid} - %msg%n"
  }
}

logger( "org.semanticweb.fbench", INFO , ["PROCFLOW"], false )

