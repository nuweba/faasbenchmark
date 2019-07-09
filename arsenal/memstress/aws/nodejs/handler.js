const MEGABYTE = 1024 * 1024;

function memIntensiveCalculation(level) {
    let amountInMB = level * 64 - 1;
    Buffer.alloc(amountInMB * MEGABYTE);
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

    memIntensiveCalculation(intensityLevel);

    let retval = {
        "reused": isWarm(),
    };

    let end = process.hrtime(startTime);
    retval.duration = end[1] + (end[0] * 1e9);
    return retval;
};


