appender( "PROCFLOW", ConsoleAppender ) {
  encoder( PatternLayoutEncoder) {
    pattern = "%.-1level - %date{ISO8601} - [%10.10thread] - %20.-20logger{0} - %.12X{uuid} - %msg%n"
  }
}

logger( "org.semanticweb.fbench", INFO , ["PROCFLOW"], false )

