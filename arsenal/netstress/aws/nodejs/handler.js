const http = require('http');

async function networkIntensive(baseNumber) {
    var startTime = process.hrtime();

    var iterationCount = 500 * Math.pow(baseNumber, 3);
    for (var i = iterationCount; i >= 0; i--) {
        await new Promise((resolve, reject) => http.get({hostname: 'google.com'}, (res) => {
            resolve(res);
        }));
    }
    var end = process.hrtime(startTime);
    return end[1] + (end[0] * 1e9);
}

function isWarm() {
    var is_warm = process.env.warm ? true : false;
    process.env.warm = true;
    return is_warm;
}

exports.handler = async (event) => {
    let intensityLevel = parseInt(event["level"]);

    return {
        "reused": isWarm(),
        "duration": await networkIntensive(intensityLevel)
    };
};
