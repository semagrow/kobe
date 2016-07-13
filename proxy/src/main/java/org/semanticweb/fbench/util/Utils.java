package org.semanticweb.fbench.util;

/**
 * Created by angel on 11/7/2016.
 */
public class Utils {

    public static long getPID() {
        String processName = java.lang.management.ManagementFactory.getRuntimeMXBean().getName();
        return Long.parseLong(processName.split("@")[0]);
    }
}
