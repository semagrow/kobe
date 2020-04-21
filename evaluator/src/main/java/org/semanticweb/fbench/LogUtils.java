package org.semanticweb.fbench;

import org.slf4j.MDC;

import java.sql.Timestamp;
import java.util.UUID;

public final class LogUtils {

    public static String getNewQueryID() {
        MDC.put("qid", UUID.randomUUID().toString().substring(0,8));
        return MDC.get("qid");
    }

    public static String getQueryID() {
        return MDC.get("qid");
    }

    public static String getCurrTime() {
        Timestamp ts = new Timestamp(System.currentTimeMillis());
        String ts_str = ts.toString();
        return ts_str.substring(ts_str.indexOf(' ') + 1);
    }
}
