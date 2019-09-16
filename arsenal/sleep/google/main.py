import json
import os
import time
import datetime

def get_current_epoch():
    return int((datetime.datetime.utcnow() - datetime.datetime(1970, 1, 1)).total_seconds() * 1000)

def get_sleep_parameter(request):
    user_input = str(request.args.get("sleep"))
    if not user_input or not user_input.isdigit() or int(user_input) < 0:
        return json.dumps({"error": "invalid sleep parameter"})
    return int(user_input) 

def get_parameters(request):
    return get_sleep_parameter(request)

def run_test(sleep_time):
    time.sleep(sleep_time / 1000.0)

def is_warm():
    is_warm = os.environ.get("warm") == "true"
    os.environ["warm"] = "true"
    return is_warm


def http(request):
    start = get_current_epoch()
    reused = is_warm()
    sleep_time = get_parameters(request)
    if type(sleep_time) != int:
        return sleep_time
    run_test(sleep_time)
    duration = (get_current_epoch() - start) * 1000000

    return json.dumps({
        "duration": duration,
        "reused": reused
    })

main128  = http
main256  = http
main512  = http
main1024 = http
main2048 = http
