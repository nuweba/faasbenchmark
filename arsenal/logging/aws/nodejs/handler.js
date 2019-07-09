function logging(baseNumber) {
    var iterationCount = 30000 * Math.pow(baseNumber, 3);
    for (var i = iterationCount; i >= 0; i--) {
        console.log('this is a log message');
    }
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    var startTime = process.hrtime();
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if(!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"}
    }

    logging(intensityLevel);

    let retval = {
        "reused": isWarm(),
    };

    var end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval
};

