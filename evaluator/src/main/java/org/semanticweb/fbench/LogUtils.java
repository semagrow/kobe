package org.semanticweb.fbench;

import org.semanticweb.fbench.query.Query;
import org.slf4j.MDC;

import java.util.UUID;

public final class LogUtils {

    public static void setMDC() {
        MDC.put("uuid", UUID.randomUUID().toString());
    }

    public static String annotateExperimentQuery(Query query, String experimentName, int run) {
        String annotation = "#kobeQueryDesc " +
                "Experiment: " + experimentName + " - " +
                "Date: " + Config.getConfig().getDate() + " - " +
                "Query: " + query.getIdentifier() + " - " +
                "Run: " + run + "\n";
        return annotation + query.getQuery();
    }
}
