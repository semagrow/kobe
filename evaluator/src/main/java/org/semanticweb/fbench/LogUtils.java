package org.semanticweb.fbench;

import org.slf4j.MDC;

import java.util.UUID;

public final class LogUtils {

    public static void setMDC() {
        MDC.put("uuid", UUID.randomUUID().toString());
    }
}
