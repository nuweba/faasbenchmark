package com.microsoft.azure.functions;

import java.util.Optional;
import com.microsoft.azure.functions.annotation.*;
import com.microsoft.azure.functions.ExecutionContext;

public class Http {
    public int getSleepParameter(HttpRequestMessage<String> req){
        String user_input = req.getQueryParameters().getOrDefault("sleep", "");
	if (user_input != "") {
		int retval = Integer.parseInt(user_input);
		return retval > -1 ? retval : -1;
	}
	return -1;
    }

    public int getParameters(HttpRequestMessage<String> req){
	return getSleepParameter(req);
    }

    public void runTest(int sleepTime) {
	try {
		Thread.sleep(sleepTime);
	}
	catch(Exception e) {}
    }

    public Boolean isWarm() {
	Boolean result = System.getProperty("warm") == "true";
	System.setProperty("warm", "true");
	return result;
    }

    @FunctionName("JavaLinuxSleep")
    public String hello(@HttpTrigger(name = "req", methods = {"post"}, authLevel = AuthorizationLevel.ANONYMOUS) HttpRequestMessage<String> req,
                        ExecutionContext context) {
        long start = System.currentTimeMillis();
	int sleepTime = getParameters(req);
	if(sleepTime < 0)
		return "{\"error\": \"invalid sleep parameter\"}";
	runTest(sleepTime);
	Boolean reused = isWarm();
	long duration = (System.currentTimeMillis() - start) * 1000000;

	final String jsonDocument = "{\"duration\":" + duration + ", " + "\"reused\": " + reused + "}";
        return jsonDocument;
    }
}
