using System.Collections.Generic;
using System.Threading;
using Amazon.Lambda.Core;
using System;

[assembly:LambdaSerializer(typeof(Amazon.Lambda.Serialization.Json.JsonSerializer))]

namespace AwsDotnetCsharp
{
    public class Handler
    {
       public int GetSleepParameter(Dictionary<String, String> input){
          String user_input = input["sleep"];
          if (user_input != null) {
             int retval = Int32.Parse(user_input);
	     return retval > -1 ? retval : -1;
	  }
	  return -1;
       }

       public int GetParameters(Dictionary<String, String> input){
            return GetSleepParameter(input);
       }

       public void RunTest(int sleepTime) {
	    Thread.Sleep(sleepTime);
       }

       public bool isWarm() {
            bool result = System.Environment.GetEnvironmentVariable("warm") == "true";
            System.Environment.SetEnvironmentVariable("warm", "true");
            return result;
       }


       public Response Hello(Dictionary<string, string> request)
       {
	   long start = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();
	   int sleepTime = GetParameters(request);
	   RunTest(sleepTime);
	   bool reused = isWarm();
	   long end = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();
	   long duration = (end - start)*1000000;
           return new Response(duration, reused);
       }
    }

    public class Response
    {
       public long duration {get; set;}
       public bool reused {get; set;}

       public Response(long duration, bool reused){
	   this.reused = reused;
	   this.duration = duration;
       }
    }
}
