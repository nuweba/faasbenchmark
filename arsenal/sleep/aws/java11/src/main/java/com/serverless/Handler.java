package com.serverless;

import java.util.Collections;
import java.util.Map;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;

class Result {
	public long duration;
	public Boolean reused;
}

public class Handler implements RequestHandler<Map<String, String>, Result> {
	public int getSleepParameter(Map<String, String> input){
		String user_input = input.get("sleep");
		if (user_input != null) {
			int retval = Integer.parseInt(user_input);
			return retval > -1 ? retval : -1;
		}
		return -1;
	}

	public int getParameters(Map<String, String> input){
		return getSleepParameter(input);
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

	@Override
	public Result handleRequest(Map<String, String> input, Context context) {
		Result retval = new Result();
		long start = System.currentTimeMillis();
		int sleepTime = getParameters(input);
		runTest(sleepTime);
		retval.reused = isWarm();
		retval.duration = (System.currentTimeMillis() - start) * 1000000;

		return retval;
	}
}
