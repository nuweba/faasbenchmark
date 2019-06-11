const MEGABYTE = 1024 * 1024;

function memIntensiveCalculation(level) {
    var startTime = process.hrtime();
    let amountInMB = level * 64 - 1;
    Buffer.allocUnsafe(amountInMB * MEGABYTE);
    let end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    let intensityLevel = event.level ? parseInt(event.level) : null;
    if(!intensityLevel || intensityLevel < 1) {
        return {"error": "invalid level parameter"}
    }

    return {
        "reused": isWarm(),
        "duration": memIntensiveCalculation(intensityLevel)
    };
};


