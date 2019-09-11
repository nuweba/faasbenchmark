using System;
using System.Net;
using System.Net.Http;
using System.Text;
using System.IO;
using System.Threading.Tasks;
using System.Threading;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Azure.WebJobs;
using Microsoft.Azure.WebJobs.Extensions.Http;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Logging;
using Newtonsoft.Json;

namespace myfunc
{
    public static class sleep
    {
       public static int GetSleepParameter(IQueryCollection input){
          string user_input = input["sleep"];
          if (user_input != null) {
             int retval = Int32.Parse(user_input);
	     return retval >= 0 ? retval : -1;
	  }
	  return -1;
       }

       public static int GetParameters(IQueryCollection input){
	       return GetSleepParameter(input);
       }

       public static void RunTest(int sleepTime) {
	       Thread.Sleep(sleepTime);
       }

       public static bool isWarm() {
	       bool result = System.Environment.GetEnvironmentVariable("warm") == "true";
	       System.Environment.SetEnvironmentVariable("warm", "true");
	       return result;
       }


        [FunctionName("CsharpLinuxSleep")]
        public static HttpResponseMessage Run(
            [HttpTrigger(AuthorizationLevel.Anonymous, "get", "post", Route = null)] HttpRequest req,
            ILogger log)
        {
	   long start = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();
	   string jsonToReturn;
	   int sleepTime = GetParameters(req.Query);
	   if (sleepTime < 0) {
		   jsonToReturn = JsonConvert.SerializeObject(new {error = "invalid sleep parameter"});
		   return new HttpResponseMessage(HttpStatusCode.OK) {
			   Content = new StringContent(jsonToReturn, Encoding.UTF8, "application/json")
		   };
	   }
	   RunTest(sleepTime);
	   bool reused = isWarm();
	   long end = new DateTimeOffset(DateTime.UtcNow).ToUnixTimeMilliseconds();
	   long duration = (end - start)*1000000;
	   jsonToReturn = JsonConvert.SerializeObject(new {duration = duration, reused = reused});
	   return new HttpResponseMessage(HttpStatusCode.OK) {
		   Content = new StringContent(jsonToReturn, Encoding.UTF8, "application/json")
	   };
        }
    }
}
